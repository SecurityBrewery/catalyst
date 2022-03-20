<template>
  <div>
    <v-row class="mt-4">
      <v-spacer></v-spacer>
      <v-btn href="https://catalyst-soar.com/docs/catalyst/analyst/tickets" target="_blank" outlined rounded small class="mr-2">
        <v-icon>mdi-book-open</v-icon> Handbook
      </v-btn>
    </v-row>

    <div v-if="$route.params.id === 'new'">
      <TicketNew></TicketNew>
    </div>
    <div v-else-if="ticket === undefined" class="text-sm-center py-16">
      <v-progress-circular
        indeterminate
        color="primary"
        :size="70"
        :width="7"
        class="align-center"
      >
      </v-progress-circular>
    </div>
    <div v-else>
      <!--v-row-- class="mt-5">
        <v-col>
          <div style="display: flex; justify-content: flex-end" class="pb-1">
            <v-btn
              :to="{ name: 'Graph', params: { col: 'tickets', id: ticket.id } }"
              class="mr-2"
              outlined
              elevation="0"
            >
              <v-icon class="mr-1">mdi-graph</v-icon> open graph
            </v-btn>
          </div>
        </v-col>
      </v-row-->
      <v-row class="mt-2">
        <v-col cols="12" lg="8">
          <h1 class="d-flex">
            <v-menu offset-y>
              <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on" class="mr-2">
                  {{ ticket.type | capitalize }}
                </span>
              </template>
              <v-list>
                <v-list-item
                  dense
                  link
                  v-for="state in otherTypes"
                  :key="state.id"
                  @click="setType(state.id)"
                >
                  <v-list-item-title>
                    Move to <v-icon small>{{ typeIcon(state.id) }}</v-icon>
                    {{ state.name }}
                  </v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
            #{{ ticket.id }}:
            <span v-if="!editName" @click="editName = ticket.name" class="ml-2" style="cursor: pointer">
              {{ ticket.name }}
            </span>
            <v-text-field
              v-else
              hide-details
              outlined
              dense
              v-model="editName"
              @blur="editName = ''"
              autofocus
              @click:append.prevent="saveName"
              @keydown.enter="saveName"
              append-icon="mdi-content-save"
              class="ml-2"
            ></v-text-field>
          </h1>

          <v-divider class="my-2"></v-divider>

          <h2 class="text--disabled" style="font-size: 12pt">
            <v-menu offset-y class="mr-2">
              <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on">
                  <v-icon small class="mr-1" :color="statusColor(ticket.status)">
                    {{ statusIcon(ticket.status) }}
                  </v-icon>
                  <span :class="statusColor(ticket.status) + '--text'">
                    {{ ticket.status | capitalize }}
                  </span>
                </span>
              </template>
              <v-list>
                <v-list-item dense link v-for="state in otherStates" :key="state" @click="setStatus(state)">
                  <v-list-item-title>
                    Set status to
                    <v-icon small>{{ statusIcon(state) }}</v-icon>
                    {{ state | capitalize }}
                  </v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
            &middot;
            <v-icon small class="mr-1">mdi-calendar-plus</v-icon>
            {{ ticket.created | formatdate($store.state.settings.timeformat) }}
            &middot;
            <v-icon small class="mr-1">mdi-calendar-edit</v-icon>
            {{ ticket.modified | formatdate($store.state.settings.timeformat) }}
          </h2>

          <v-divider class="mt-0 mb-4"></v-divider>

          <v-dialog v-model="schemaDialog"  max-width="600px">
            <v-card>
              <v-card-title>Select new Template</v-card-title>
              <v-card-text>
                <v-select
                    label="Template"
                    :items="templates"
                    item-text="name"
                    return-object
                    v-model="selectedTemplate"
                ></v-select>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn class="mr-1" small dense color="error" outlined @click="schemaDialog = false">
                    <v-icon>mdi-close</v-icon>
                    Cancel
                  </v-btn>
                  <v-btn small dense color="success" outlined @click="setSchema(selectedTemplate)">
                    <v-icon>mdi-content-save</v-icon>
                    Save
                  </v-btn>
                </v-card-actions>
              </v-card-text>
            </v-card>
          </v-dialog>
          <div style="align-items: center" class="d-flex pb-1">
            <span class="text--disabled">Details</span>
            <v-spacer></v-spacer>
            <v-btn small text @click="schemaDialog = true">
              Change Template
            </v-btn>
          </div>

          <v-form v-model="valid" class="clearfix pt-2 pb-0 mb-0">
            <v-jsf
              v-model="ticket.details"
              :schema="schema"
              :options="{ initialValidation: 'all', formats: { time: timeformat, date: dateformat, 'date-time': datetimeformat }, editMode: 'inline' }"
            />
          </v-form>
          <v-btn small class="float-right mb-2" color="card" @click="saveTicket" outlined>
            <v-icon class="mr-1">mdi-content-save</v-icon>
            Save Details
          </v-btn>

          <v-divider class="my-2 mb-4" style="clear:both;"></v-divider>

          <span class="text--disabled">Log</span>

          <div>
            <v-textarea
                v-model="comment"
                hide-details
                flat
                label="Add a comment..."
                solo
                auto-grow
                rows="2"
                class="py-2"
            >
              <template v-slot:append>
                <v-btn class="mx-0 mt-n1" text @click="addComment">
                  <v-icon>mdi-send</v-icon>
                </v-btn>
              </template>
            </v-textarea>
            <div v-for="(comment, id) in logs(ticket)" :key="id" class="pb-2">
              <div v-if="'type' in comment && comment.message !== 'AddComment'" style="text-align: center">
                <span class="text--disabled" :title="comment.created">
                  {{ comment.message }} &middot;
                  <strong> {{ comment.creator }}</strong> &middot;
                  {{ relDate(comment.created) }}
                </span>
              </div>
              <v-card v-else-if="!('type' in comment)" elevation="0" color="cards">
                <v-card-subtitle class="pb-0">
                  <strong> {{ comment.creator }}</strong>
                  <span class="text--disabled ml-3" :title="comment.created">
                  {{ relDate(comment.created) }}
                </span>
                </v-card-subtitle>
                <v-card-text class="mb-0 mt-2">
                  <vue-markdown>{{ replacelinks(comment.message) }}</vue-markdown>
                </v-card-text>
              </v-card>
            </div>
          </div>
        </v-col>
        <v-col cols="12" lg="4" >
          <v-btn
            v-if="$store.state.settings.tier === 'enterprise'"
            :to="{ name: 'Graph', params: { col: 'tickets', id: ticket.id } }"
            class="mr-2 mb-2"
            outlined
            small
            elevation="0"
          >
            <v-icon class="mr-1">mdi-graph</v-icon> open graph
          </v-btn>

          <v-list dense color="background">
            <v-list-item class="pa-0 ma-0" style="min-height: 32px">
              <span class="text--disabled" style="width: 50px;">Owner</span>
              <v-menu offset-y class="mr-2">
                <template v-slot:activator="{ on, attrs }">
                  <span v-bind="attrs" v-on="on">
                    <User :id="ticket.owner" class="ml-3"></User>
                  </span>
                </template>
                <v-list>
                  <v-list-item
                      dense
                      link
                      v-for="user in otherUsers(ticket.owner)"
                      :key="user.id"
                      @click="setOwner(user.id)"
                  >
                    <v-list-item-title>
                      Change owner to <User :id="user.id"></User>
                    </v-list-item-title>
                  </v-list-item>
                  <v-list-item
                      v-if="ticket.owner"
                      dense
                      link
                      @click="setOwner(undefined)"
                  >
                    <v-list-item-title>
                      Unassign ticket
                    </v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </v-list-item>
            <!--v-list-item class="pa-0 ma-0" style="min-height: 32px">
              <span class="text--disabled" style="width: 50px;">Editors</span>
              <span v-for="writer in ticket.write" :key="writer">
                <User :id="writer" class="ml-3"></User>
              </span>
              <v-btn v-if="!ticket.write || ticket.write.length === 0" small text elevation="0" @click="referenceDialog = true">
                <v-icon small>mdi-plus</v-icon> Add editor
              </v-btn>
            </v-list-item>
            <v-list-item class="pa-0 ma-0" style="min-height: 32px">
              <span class="text--disabled" style="width: 50px;">Viewers</span>
              <span v-for="reader in ticket.read" :key="reader">
                <User :id="reader" class="ml-3"></User>
              </span>
              <v-btn v-if="!ticket.read || ticket.read.length === 0" small text elevation="0" @click="referenceDialog = true">
                <v-icon small>mdi-plus</v-icon> Add viewer
              </v-btn>
            </v-list-item>
          </v-list-->

          <v-divider class="mb-5"></v-divider>

          <div style="align-items: center" class="d-flex pb-1">
            <span class="text--disabled">Playbooks</span>
            <v-spacer></v-spacer>
            <v-btn small icon @click="playbookDialog = true">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </div>
          <div
            v-for="(playbook, playbookid) in ticket.playbooks"
            :key="playbookid"
            class="mb-2"
          >
            <v-card color="cards" flat>
              <v-card-subtitle class="d-flex pa-2 pb-1 mb-0" style="line-height: 28px;">
                {{ playbook.name }}
                <v-spacer></v-spacer>
                <v-btn
                  small
                  icon
                  @click="
                    deletePlaybookID = playbookid;
                    deletePlaybookDialog = true;
                  "
                >
                  <v-icon small>mdi-close-circle</v-icon>
                </v-btn>
              </v-card-subtitle>
              <div style="overflow-x: scroll">
                <vue-pipeline
                  v-if="showPipelines"
                  ref="pipeline"
                  :x="50"
                  :y="55"
                  :data="pipeline(playbook)"
                  :showArrow="true"
                  :ystep="70"
                  :xstep="90"
                  lineStyle="default"
                  @select="select"
                  class="mx-4"
                />
              </div>
              <v-list dense color="cards" class="py-0">
                <v-list-item
                  v-for="taskwithid in stasks(playbookid)"
                  :key="taskwithid.id"
                  color="blue"
                  no-action
                  link
                  @click="selectTask(taskwithid, playbookid)"
                >
                  <!--template v-slot:activator-->
                  <v-list-item-icon>
                    <v-icon
                      :class="{
                        spinning: taskwithid.task.active && taskwithid.task.type === 'automation'
                      }"
                      :disabled="!taskwithid.task.active"
                      v-text="icon(taskwithid.task)"
                    />
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title
                      :class="{
                        strike: taskwithid.task.done, 'text--disabled': !taskwithid.task.active
                      }"
                    >
                      {{ taskwithid.task.name }}

                      <span v-if="taskwithid.task.owner" class="text--disabled float-right">
                        {{ taskwithid.task.owner }}
                      </span>
                    </v-list-item-title>
                    <v-list-item-subtitle>
                      <div v-if="(taskwithid.task.data && 'error' in taskwithid.task.data)">
                        <b>Error:</b> {{ taskwithid.task.data.error }}
                      </div>
                    </v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
              </v-list>
            </v-card>
          </div>
          <v-dialog v-model="taskDialog" max-width="600px">
            <v-card v-if="selectedTask">
              <v-card-title>
                <span class="headline">{{ selectedTask.task.name }}</span>
              </v-card-title>
              <v-card-subtitle>
                <v-menu offset-y class="mr-2">
                  <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on">
                  <v-icon small class="mr-1">mdi-account</v-icon>
                  {{ selectedTask.task.owner ? selectedTask.task.owner : "unassigned" }}
                </span>
                  </template>
                  <v-list>
                    <v-list-item
                        dense
                        link
                        v-for="user in otherUsers(selectedTask.task.owner)"
                        :key="user.username"
                        @click="setTaskOwner(selectedTaskPlaybook, selectedTask.id, user.id)"
                    >
                      <v-list-item-title>
                        Change owner to {{ user.id }}
                      </v-list-item-title>
                    </v-list-item>
                    <v-list-item
                        v-if="selectedTask.task.owner"
                        dense
                        link
                        @click="setTaskOwner(selectedTaskPlaybook, selectedTask.id, '')"
                    >
                      <v-list-item-title>
                        Unassign task
                      </v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
                &middot;
                <v-icon small class="mr-1">mdi-calendar-plus</v-icon>
                {{ selectedTask.task.created | formatdate($store.state.settings.timeformat) }}
                <span v-if="selectedTask.task.closed">
                  &middot;
                  <v-icon small class="mr-1">mdi-calendar-check</v-icon>
                  {{ selectedTask.task.closed | formatdate($store.state.settings.timeformat) }}
                </span>
              </v-card-subtitle>

              <v-card-text>
                <v-form
                    v-if="(selectedTask.task.type === 'input' || selectedTask.task.type === 'task')"
                    :ref="selectedTask.id"
                    v-model="valids[selectedTaskPlaybook.toString() + '-' + selectedTask.id]"
                    class="clearfix"
                >
                  <h3 v-if="selectedTask.task.schema">Input</h3>
                  <v-jsf
                      v-if="selectedTask.task.schema"
                      v-model="tdata[selectedTaskPlaybook.toString() + '-' + selectedTask.id]"
                      :schema="selectedTask.task.schema"
                  />
                </v-form>
                <div v-if="(selectedTask.task.data && 'error' in selectedTask.task.data)">
                  <h3>Error</h3>
                  <pre>{{ selectedTask.task.data.error }}</pre>
                </div>
                <div v-else-if="selectedTask.task.data">
                  <h3>Data</h3>
                  <JSONHTML :json="selectedTask.task.data"></JSONHTML>
                </div>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn class="float-right mr-1" small dense color="error" outlined @click="taskDialog = false">
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
                <v-btn
                    class="float-right"
                    v-if="selectedTaskPlaybook.toString() + '-' + selectedTask.id in valids &&
                            valids[selectedTaskPlaybook.toString() + '-' + selectedTask.id]"
                    small
                    dense
                    color="success"
                    outlined
                    @click="complete(selectedTaskPlaybook, selectedTask.id)"
                >
                  <v-icon>mdi-check</v-icon>
                  Complete
                </v-btn>
                <v-btn class="float-right" v-else small dense color="info" outlined @click="save(selectedTaskPlaybook, selectedTask.id)">
                  <v-icon>mdi-content-save</v-icon>
                  Save
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="playbookDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline">Add Playbook</span>
              </v-card-title>
              <v-card-text>
                <v-select
                  :items="playbooks"
                  item-text="name"
                  item-value="id"
                  return-object
                  filled
                  label="Playbook"
                  v-model="selectedPlaybook"
                ></v-select>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" outlined @click="playbookDialog = false">
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
                <v-btn color="success" outlined @click="addPlaybook">
                  <v-icon>mdi-content-save</v-icon>
                  Save
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="deletePlaybookDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline">
                  Remove Playbook {{ deletePlaybookID }}
                </span>
              </v-card-title>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="error"
                  outlined
                  @click="deletePlaybookDialog = false"
                >
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
                <v-btn color="success" outlined @click="removePlaybook">
                  <v-icon>mdi-delete</v-icon>
                  Remove
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-divider class="my-5"></v-divider>

          <div style="display: flex; align-items: center" class="pb-1">
            <span class="text--disabled">References</span>
            <v-spacer></v-spacer>
            <v-btn small icon @click="referenceDialog = true">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </div>
          <v-list dense color="background">
            <v-list-item
              class="pa-0 ma-0"
              style="min-height: 22px"
              v-for="reference in ticket.references"
              :key="reference.href"
            >
              <v-list-item-title :title="reference.href">
                <a class="text--primary" :href="reference.href" target="_blank">
                  <span v-if="reference.name">
                    {{ reference.name }}
                    <span style="opacity: 0.3" class="ml-2">
                      {{ reference.href }}
                    </span>
                  </span>
                  <span v-else>{{ reference.href }}</span>
                </a>
              </v-list-item-title>
              <v-list-item-action class="ma-0">
                <v-btn icon x-small>
                  <v-icon
                    small
                    @click="
                      deleteReferenceID = reference.href;
                      deleteReferenceDialog = true;
                    "
                  >
                    mdi-close-circle
                  </v-icon>
                </v-btn>
              </v-list-item-action>
            </v-list-item>
          </v-list>
          <v-dialog v-model="referenceDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline">Add Reference</span>
              </v-card-title>
              <v-card-text>
                <v-text-field
                  label="Title"
                  v-model="newName"
                  required
                  placeholder="e.g. My External Page"
                ></v-text-field>
                <v-text-field
                  label="Link"
                  v-model="newLink"
                  required
                  placeholder="e.g. https://example.org/asd"
                  @keydown.enter="addReference"
                ></v-text-field>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" outlined @click="referenceDialog = false">
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
                <v-btn color="success" outlined @click="addReference">
                  Save
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="deleteReferenceDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline"
                  >Remove Reference {{ deleteReferenceID }}</span
                >
              </v-card-title>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="error"
                  outlined
                  @click="deleteReferenceDialog = false"
                >
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
                <v-btn color="success" outlined @click="removeReference">
                  <v-icon>mdi-delete</v-icon>
                  Remove
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-divider class="my-2 mb-4"></v-divider>

          <div style="display: flex; align-items: center" class="pb-1">
            <span class="text--disabled">Artifacts</span>
            <v-spacer></v-spacer>
            <v-btn small icon @click="artifactDialog = true">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </div>
          <v-list dense v-if="ticket.artifacts && ticket.artifacts.length">
            <ArtifactSnippet
              v-for="artifact in ticket.artifacts"
              :key="artifact.name"
              xto="{ name: 'ArtifactPopup', params: { type: '-', id: ticket.id, artifact: artifact.name } }"
              @click="showArtifactDialog = true; selectedArtifact = artifact.name"
              :artifact="artifact"
              action="mdi-close-circle"
              @actionClick="
                deleteArtifactID = artifact.name;
                deleteArtifactDialog = true;
              "
            ></ArtifactSnippet>
          </v-list>
          <v-dialog v-model="artifactDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline">Add Artifact</span>
              </v-card-title>
              <v-card-text>
                <v-text-field
                  label="Artifact"
                  v-model="newArtifact"
                  required
                  placeholder="e.g. 8.8.8.8, example.org, 3e63b753e36e74622a4dcaac34cfc613"
                  @keydown.enter="addArtifact"
                ></v-text-field>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" outlined @click="artifactDialog = false">
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
                <v-btn color="success" outlined @click="addArtifact">
                  <v-icon>mdi-content-save</v-icon>
                  Save
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="deleteArtifactDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline"
                  >Remove Artifact Link to {{ deleteArtifactID }}</span
                >
              </v-card-title>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="error"
                  outlined
                  @click="deleteArtifactDialog = false"
                >
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
                <v-btn color="success" outlined @click="removeArtifact">
                  <v-icon>mdi-delete</v-icon>
                  Remove
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="showArtifactDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline">Artifact</span>
              </v-card-title>
              <v-card-text>
                <ArtifactPopup :id="ticket.id" :name="selectedArtifact"></ArtifactPopup>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                    color="error"
                    outlined
                    @click="showArtifactDialog = false"
                >
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-divider class="my-2 mb-4"></v-divider>

          <div style="display: flex; align-items: center" class="pb-1">
            <span class="text--disabled">Related Tickets</span>
            <v-spacer></v-spacer>
            <v-btn small icon @click="relatedDialog = true">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </div>
          <v-list dense v-if="ticket.tickets && ticket.tickets.length">
            <TicketSnippet
              v-for="relatedTicket in ticket.tickets"
              :key="relatedTicket.id"
              :to="{ name: 'Ticket', params: { id: relatedTicket.id } }"
              :ticket="relatedTicket"
              action="mdi-close-circle"
              @actionClick="
                deleteRelatedID = relatedTicket.id;
                deleteRelatedDialog = true;
              "
            ></TicketSnippet>
          </v-list>
          <v-dialog v-model="relatedDialog" max-width="800px">
            <v-card>
              <v-card-title>
                <span class="headline">Add Ticket</span>
              </v-card-title>
              <v-card-text>
                <TicketList @click="linkTicket"></TicketList>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" outlined @click="relatedDialog = false">
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="deleteRelatedDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline"
                  >Remove Ticket Relation to {{ deleteRelatedID }}</span
                >
              </v-card-title>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="error"
                  outlined
                  @click="deleteRelatedDialog = false"
                >
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
                <v-btn color="success" outlined @click="removeRelation">
                  <v-icon>mdi-delete</v-icon>
                  Remove
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-divider class="my-2 mb-4"></v-divider>

          <div style="display: flex; align-items: center" class="pb-1">
            <span class="text--disabled">Files</span>
            <v-spacer></v-spacer>
            <v-btn v-if="hasRole('analyst:file')" small icon @click="fileDialog = true">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </div>
          <v-list dense color="background">
            <v-list-item
              v-for="file in ticket.files"
              :key="file.key"
              class="pa-0 ma-0"
              style="min-height: 20px"
              autofocus
            >
              <v-list-item-title>
                {{ file.name }}
              </v-list-item-title>
              <v-list-item-action class="ma-0">
                <v-btn
                  icon
                  x-small
                  @click="download(file.name, file.key)"
                  style="opacity: 0"
                >
                  <v-icon>mdi-download</v-icon>
                </v-btn>
              </v-list-item-action>
            </v-list-item>
          </v-list>
          <v-dialog v-model="fileDialog" max-width="600px">
            <v-card>
              <v-card-title>
                <span class="headline">Add File</span>
              </v-card-title>
              <v-card-text>
                <dashboard
                  :uppy="uppy"
                  :props="{ theme: 'dark', height: 400 }"
                />
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="red" outlined @click="fileDialog = false">
                  <v-icon>mdi-close</v-icon>
                  Cancel
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-col>
      </v-row>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import {alg, Graph} from "graphlib";

import ArtifactSnippet from "../components/snippets/ArtifactSnippet.vue";
import TicketSnippet from "../components/snippets/TicketSnippet.vue";
import TicketList from "../components/TicketList.vue";
import User from "../components/User.vue";
import ArtifactPopup from "./ArtifactPopup.vue";

import {
  Ticket,
  TicketResponse,
  TicketTemplate,
  PlaybookTemplate,
  Reference,
  Task,
  Type,
  TypeColorEnum,
  TaskResponse, PlaybookResponse, UserResponse, TicketWithTickets,
} from "@/client";
import {API} from "@/services/api";

import {Dashboard} from "@uppy/vue";
import Tus from "@uppy/tus";

import "@uppy/core/dist/style.css";
import "@uppy/dashboard/dist/style.css";

import {Uppy} from "@uppy/core";
import {AxiosError, AxiosResponse} from "axios";
import {DateTime} from "luxon";
import VueMarkdown from "vue-markdown";
import JSONHTML from "../components/JSONHTML.vue";
import TicketNew from "@/views/TicketNew.vue";

interface State {
  valid: boolean;

  editName: string;
  readOnly: boolean;
  ticket?: TicketResponse;

  templates: Array<TicketTemplate>;
  selectedTemplate?: TicketTemplate;

  schemaDialog: false;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  g: Record<string, any>;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  selected: any;

  showPipelines: boolean;
  playbookDialog: boolean;
  deletePlaybookDialog: boolean;
  deletePlaybookID: string;
  playbooks: Array<PlaybookTemplate>;
  selectedPlaybook?: PlaybookTemplate;
  selectedTask?: TaskWithID;
  selectedTaskPlaybook?: string;
  taskDialog: boolean;

  newArtifact: string;
  artifactDialog: boolean;
  deleteArtifactDialog: boolean;
  deleteArtifactID: string;
  showArtifactDialog: boolean;
  selectedArtifact: string;

  relatedDialog: boolean;
  deleteRelatedDialog: boolean;
  deleteRelatedID: number;

  tdata: Record<string, never>;
  valids: Record<string, boolean>;

  newName: string;
  newLink: string;
  referenceDialog: boolean;
  deleteReferenceID: string;
  deleteReferenceDialog: boolean;

  newFile: string;
  fileDialog: boolean;

  uppy?: Uppy;

  users: Array<UserResponse>;

  comment: string;
}

interface TaskWithID {
  id: string;
  task: TaskResponse;
}

export default Vue.extend({
  name: "Ticket",
  components: {
    TicketNew,
    Dashboard,
    ArtifactSnippet,
    TicketSnippet,
    "vue-markdown": VueMarkdown,
    TicketList,
    ArtifactPopup,
    JSONHTML,
    User
  },
  data: (): State => ({
    valid: false,

    editName: "",
    readOnly: false,
    ticket: undefined,

    templates: [],
    selectedTemplate: undefined,

    schemaDialog: false,

    g: {},
    selected: undefined,

    showPipelines: true,
    playbookDialog: false,
    deletePlaybookDialog: false,
    deletePlaybookID: "",
    playbooks: [],
    selectedPlaybook: undefined,
    selectedTask: undefined,
    selectedTaskPlaybook: undefined,
    taskDialog: false,

    newArtifact: "",
    artifactDialog: false,
    deleteArtifactDialog: false,
    deleteArtifactID: "",
    showArtifactDialog: false,
    selectedArtifact: "",

    relatedDialog: false,
    deleteRelatedDialog: false,
    deleteRelatedID: -1,

    tdata: {},
    valids: {},

    newName: "",
    newLink: "",
    referenceDialog: false,
    deleteReferenceDialog: false,
    deleteReferenceID: "",

    newFile: "",
    fileDialog: false,

    uppy: undefined,

    users: [],

    comment: ""
  }),
  watch: {
    $route: function(): void {
      this.loadTicket();
    }
  },
  computed: {
    schema: function() {
      if (this.ticket !== undefined && this.ticket.schema !== undefined) {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        let schema: any = JSON.parse(this.ticket.schema);
        schema.readOnly = this.readOnly;
        return schema;
      }
      return {};
    },
    otherStates: function(): Array<string> {
      if (!this.ticket || !this.ticket.status) {
        return [];
      }
      if (this.ticket.status === "closed") {
        return ["open"];
      }
      return ["closed"];
    },
    otherTypes: function(): Array<Type> {
      return this.lodash.filter(
        this.$store.state.settings.ticketTypes,
        (state: Type) => {
          if (!this.ticket || !this.ticket.type) {
            return true;
          }
          return state.id !== this.ticket.type;
        }
      );
    }
  },
  methods: {
    hasRole: function (s: string) {
      if (this.$store.state.user.roles) {
        return this.lodash.includes(this.$store.state.user.roles, s);
      }
      return false;
    },
    otherUsers: function(owner: string): Array<UserResponse> {
      return this.lodash.filter(this.users, (user: UserResponse) => {
        return user.id !== owner && !user.apikey;
      });
    },
    addComment: function () {
      if (!this.ticket) {
        return;
      }

      API.addComment(this.ticket.id, { message: this.comment }).then((response) => {
        this.$store.dispatch("alertSuccess", { name: "Comment added" });
        this.setTicket(response.data);
      });
    },
    setSchema: function (template: TicketTemplate) {
      if (!this.ticket) {
        return;
      }

      API.setSchema(this.ticket.id, JSON.stringify(template.schema)).then((response) => {
        this.$store.dispatch("alertSuccess", { name: "New template saved" });
        this.setTicket(response.data);
      });
      this.schemaDialog = false;
    },
    selectTask: function (task: TaskWithID, playbookID: string) {
      this.selectedTask = task;
      this.taskDialog = true;
      this.selectedTaskPlaybook = playbookID;
    },
    timeformat: function(s: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    dateformat: function(s: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    datetimeformat: function(s: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    pipeline: function(playbook: PlaybookResponse): Array<any> {
      let tasks = this.tasks(playbook);
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      let elements = [] as Array<any>;
      this.lodash.forEach(tasks, task => {
        let elem = {
          id: task.id,
          name: task.task.name,
          next: [],
          status: "unknown"
        };

        if (task.task.done) {
          elem.status = "success";
        } else if (task.task.active) {
          elem.status = "running";
        } else {
          elem.status = "unknown";
        }
        // elem.status = "failure"; TODO

        elements.push(elem);
      });

      this.lodash.forEach(tasks, (task: TaskWithID) => {
        if ("next" in task.task) {
          this.lodash.forEach(task.task.next, (condition, nextKey) => {
            let nextID = this.lodash.findIndex(elements, ["id", nextKey]);
            let stepID = this.lodash.findIndex(elements, ["id", task.id]);
            if (nextID !== -1) {
              // TODO: invalid schema
              elements[stepID].next.push({ index: nextID });
            }
          });
        }
      });

      return this.lodash.values(elements);
    },
    select: function(node) {
      this.selected = node;
    },
    tasks: function(playbook: PlaybookResponse): Array<TaskWithID> {
      if (this.g[playbook.name] === undefined) {
        return [];
      }
      let taskKeys = alg.topsort(this.g[playbook.name]);
      let tasks = [] as Array<TaskWithID>;
      for (const tasksKey in taskKeys) {
        let taskWithID = {} as TaskWithID;
        if (playbook.tasks[taskKeys[tasksKey]] === undefined) {
          continue; // TODO
        }
        taskWithID.task = playbook.tasks[taskKeys[tasksKey]];
        taskWithID.id = taskKeys[tasksKey];
        tasks.push(taskWithID);
      }
      return tasks;
    },

    stasks: function(playbookID: string): Array<TaskWithID> {
      if (!this.ticket ||!this.ticket.playbooks || !this.ticket.playbooks[playbookID]) {
        return [];
      }

      let tasks: Array<TaskWithID> = [];
      this.lodash.forEach(this.ticket.playbooks[playbookID].tasks, (task, taskID) => {
        if (task.done || task.active) {
          tasks.push({ id: taskID, task: task });
        }
      });

      return this.lodash.sortBy(tasks, [
        (a: TaskWithID) => {
          return a.task.order;
        }
      ]);
    },
    complete(playbookID: string, taskID: string) {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      let data = {}
      if (playbookID.toString() + "-" + taskID in this.tdata) {
        data = this.tdata[playbookID.toString() + "-" + taskID];
      }

      API.completeTask(this.ticket.id, playbookID, taskID, data).then((response) => {
        this.$store.dispatch("alertSuccess", { name: "Task completed" });
        this.setTicket(response.data);
      });
      this.taskDialog = false;
      this.selectedTask = undefined;
      this.selectedTaskPlaybook = undefined;
    },
    save(playbookID: string, taskID: string) {
      if (!this.ticket || !this.ticket.id || !this.ticket.playbooks) {
        return;
      }

      let data = this.tdata[playbookID.toString() + "-" + taskID];
      API.setTaskData(this.ticket.id, playbookID, taskID, data).then((response) => {
        this.$store.dispatch("alertSuccess", { name: "Task saved" });
        this.setTicket(response.data);
      });
      this.taskDialog = false;
      this.selectedTask = undefined;
      this.selectedTaskPlaybook = undefined;
    },
    setTaskOwner: function (playbookID: string, taskID: string, owner: string) {
      if (!this.ticket || !this.ticket.id || !this.ticket.playbooks) {
        return;
      }

      API.setTaskOwner(this.ticket.id, playbookID, taskID, owner).then((response) => {
        this.$store.dispatch("alertSuccess", { name: "Owner saved" });
        this.setTicket(response.data);
        if (response.data.playbooks) {
          this.selectTask({id: taskID, task: response.data.playbooks[playbookID].tasks[taskID]}, playbookID);
        }
      });
    },
    gstatus: function(task: TaskResponse) {
      if (task.active) {
        return "open"
      }
      return "inactive"
    },
    icon(task: Task) {
      if (task.data && "error" in task.data) {
        return "mdi-alert-octagon";
      }
      if (!task.done) {
        if (task.type === "automation") {
          return "mdi-cog";
        }
        return "mdi-circle-outline";
      }
      if (task.type === "automation") {
        return "mdi-code-tags-check";
      }
      return "mdi-check-bold";
    },
    toTicket(ticketResponse: TicketResponse): Ticket {
      return {
        artifacts: ticketResponse.artifacts,
        comments: ticketResponse.comments,
        created: ticketResponse.created,
        details: ticketResponse.details,
        files: ticketResponse.files,
        modified: ticketResponse.modified,
        name: ticketResponse.name,
        owner: ticketResponse.owner,
        playbooks: ticketResponse.playbooks,
        read: ticketResponse.read,
        references: ticketResponse.references,
        schema: ticketResponse.schema,
        status: ticketResponse.status,
        type: ticketResponse.type,
        write: ticketResponse.write,
      } as Ticket
    },
    saveTicket() {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      API.updateTicket(this.ticket.id, this.toTicket(this.ticket)).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Ticket saved" });
        this.setTicket(response.data);
      });
    },
    setStatus(status: string) {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      let e: Ticket = this.toTicket(this.ticket);
      e.status = status;

      API.updateTicket(this.ticket.id, e).then(response => {
        this.$store.dispatch("alertSuccess", {
          name: "Ticket status changed",
          type: "success"
        });
        this.setTicket(response.data);
        this.$store.dispatch("fetchCount");
      });
    },
    setType(t: string) {
      if (!this.ticket || !this.ticket.id) {
        return;
      }
      this.ticket.type = t;

      API.updateTicket(this.ticket.id, this.toTicket(this.ticket)).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Ticket type changed" });
        this.setTicket(response.data);
        if (response.data.id !== undefined) {
          this.$router.push({
            name: "Ticket",
            params: { type: t, id: response.data.id.toString() }
          });
        }
      });
    },
    addPlaybook() {
      if (!this.ticket || !this.ticket.id || !this.selectedPlaybook) {
        return;
      }

      API.addTicketPlaybook(this.ticket.id, this.selectedPlaybook).then(
        response => {
          this.$store.dispatch("alertSuccess", { name: "Playbook added" });
          this.setTicket(response.data);
          this.playbookDialog = false;
        }
      );
    },
    removePlaybook() {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      API.removeTicketPlaybook(this.ticket.id, this.deletePlaybookID).then(
        response => {
          this.$store.dispatch("alertSuccess", { name: "Playbook removed" });
          this.setTicket(response.data);
          this.deletePlaybookDialog = false;
        }
      );
    },
    addArtifact() {
      if (!this.ticket || !this.ticket.id) {
        return;
      }
      API.addArtifact(this.ticket.id, {
        name: this.newArtifact,
        status: "unknown"
      }).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Artifact added" });
        this.setTicket(response.data);
        this.artifactDialog = false;
      });
    },
    removeArtifact() {
      if (!this.ticket || !this.ticket.id) {
        return;
      }
      API.removeArtifact(
        this.ticket.id,
        JSON.stringify(this.deleteArtifactID)
      ).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Artifact deleted" });
        this.setTicket(response.data);
        this.deleteArtifactDialog = false;
      });
    },
    addReference() {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      let references: Array<Reference> = [];
      if (this.ticket.references) {
        references = this.ticket.references;
      }

      references.push({ name: this.newName, href: this.newLink });
      API.setReferences(this.ticket.id, references).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Reference saved" });
        this.setTicket(response.data);
        this.referenceDialog = false;
      });
    },
    removeReference() {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      let references: Array<Reference> = [];
      if (this.ticket.references) {
        references = this.ticket.references;
      }

      references = this.lodash.filter(references, (o: Reference) => {
        return o.href !== this.deleteReferenceID;
      });
      API.setReferences(this.ticket.id, references).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Reference removed" });
        this.setTicket(response.data);
        this.deleteReferenceDialog = false;
      });
    },
    linkTicket: function(ticket: TicketResponse) {
      if (!this.ticket || !this.ticket.id || ticket.id === undefined) {
        return;
      }
      if (this.ticket.id === ticket.id) {
        this.$store.dispatch("alertError", { name: "Cannot link an ticket to itself" } );
        return;
      }

      API.linkTicket(this.ticket.id, ticket.id).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Ticket linked" });
        this.setTicket(response.data);
      });
    },
    removeRelation() {
      if (!this.ticket || !this.ticket.id || this.deleteRelatedID < 1) {
        return;
      }

      API.unlinkTicket(this.ticket.id, this.deleteRelatedID).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Ticket unlinked" });
        this.setTicket(response.data);
        this.deleteRelatedDialog = false;
      });
    },
    download(filename: string, key: string) {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      this.$http
        .get("/api/files/" + this.ticket.id + "/download/" + key, {
          responseType: "blob"
        })
        .then((response: AxiosResponse) => {
          const url = window.URL.createObjectURL(new Blob([response.data]));
          const link = document.createElement("a");
          link.href = url;
          link.setAttribute("download", filename);
          document.body.appendChild(link);
          link.click();

          this.$store.dispatch("alertSuccess", { name: "Download started" });
        })
        .catch((error: AxiosError) => {
          this.$store.dispatch("alertError", { name: "Download failed", detail: error });
        });
    },
    loadTicket(): void {
      let id = Number(this.$route.params.id);
      this.uppy = this.setupUppy(id);
      API.getTicket(id).then(response => {
        this.setTicket(response.data);
        this.showPipelines = false;
        this.$nextTick(() => {
          this.showPipelines = true;
        })
      });
    },
    setupUppy: function(id: number) {
      let uppy = new Uppy();
      uppy.use(Tus, {
        endpoint: location.protocol + '//' + location.hostname + ':'+ location.port + "/api/files/" + id.toString() + "/tusd"
      });
      return uppy;
    },
    typeIcon: function(typ: string) {
      let icon = "mdi-help";
      this.lodash.forEach(
        this.$store.state.settings.ticketTypes,
        (ticketType: Type) => {
          if (typ === ticketType.id) {
            icon = ticketType.icon;
          }
        }
      );
      return icon;
    },
    statusIcon: function(status: string) {
      if (status === 'closed') {
        return "mdi-checkbox-marked-circle-outline";
      }
      return "mdi-arrow-right-drop-circle-outline";
    },
    statusColor: function(status: string) {
      if (status === 'closed') {
        return TypeColorEnum.Success;
      }
      return TypeColorEnum.Info;
    },
    saveName: function() {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      this.ticket.name = this.editName;

      API.updateTicket(this.ticket.id, this.toTicket(this.ticket)).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Ticket name changed" });
        this.setTicket(response.data);
        this.editName = "";
      });
    },
    setOwner: function(owner?: string) {
      if (!this.ticket || !this.ticket.id) {
        return;
      }

      if (owner === undefined) {
        this.lodash.unset(this.ticket, "owner")
      } else {
        this.ticket.owner = owner;
      }

      API.updateTicket(this.ticket.id, this.toTicket(this.ticket)).then(response => {
        this.$store.dispatch("alertSuccess", { name: "Ticket owner changed" });
        this.setTicket(response.data);
      });
    },
    setTicket: function(ticket: TicketResponse) {
      this.lodash.forEach(ticket.playbooks, (playbook: PlaybookResponse) => {
        let g = new Graph();

        for (const stepKey in playbook.tasks) {
          g.setNode(stepKey);
        }

        this.lodash.forEach(playbook.tasks, (task: TaskResponse, stepKey: string) => {
          if ("next" in task) {
            this.lodash.forEach(task.next, (condition, nextKey) => {
              g.setEdge(stepKey, nextKey);
            });
          }
        });

        this.g[playbook.name] = g;
      });
      this.ticket = ticket;
    },
    relDate: function(date: string) {
      let rtf = new Intl.RelativeTimeFormat("en", { numeric: "auto" });
      let deltaDays =
          (new Date(date).getTime() - new Date().getTime()) / (1000 * 3600 * 24);
      let relDate = rtf.format(Math.round(deltaDays), "days");
      if (deltaDays > -3) {
        relDate +=
            ", " +
            new Date(date).toLocaleTimeString([], {
              hour: "2-digit",
              minute: "2-digit"
            });
      }
      return relDate;
    },
    logs: function(ticket: TicketWithTickets) {
      return this.lodash.reverse(this.lodash.sortBy(this.lodash.union(ticket.comments, ticket.logs), ['created']))
    },
    replacelinks: function (text: string): string {
      return text.replace(/#(\d+)\b/, (match, p1) => {
        return "[" + match + "](/tickets/" + p1 + ")";
      })
    }
  },
  mounted() {
    if (this.$route.params.id === 'new') {
      return
    }

    this.loadTicket();
    API.listUsers().then(response => {
      this.users = response.data;
    });
    API.listTemplates().then(response => {
      this.templates = response.data;
      if (this.templates && this.templates.length > 0) {
        this.selectedTemplate = this.templates[0];
      }
    });
    API.listPlaybooks().then(response => {
      this.playbooks = response.data;
      if (this.playbooks && this.playbooks.length > 0) {
        this.selectedPlaybook = this.playbooks[0];
      }
    });

    this.$store.subscribeAction((automation) => {
      if (!automation.payload || !("ids" in automation.payload) || !automation.payload["ids"]) {
        return
      }
      Vue.lodash.forEach(automation.payload["ids"], id => {
        if (!this.ticket || !this.ticket.id) {
          return;
        }
        if (id === "tickets/" + this.ticket.id) {
          this.loadTicket();
        }
      });
    });
  },
  beforeDestroy() {
    if (this.uppy !== undefined) {
      this.uppy.close();
    }
  }
});
</script>

<style>
.pipeline-node-label {
  width: 120px !important;
  fill: #ffffff !important;
}

.svgResultStatus > circle {
  /* stroke: rgb(48, 48, 48) !important; */
  fill: rgb(33, 33, 33);
}

.pipeline-node-selected .svgResultStatus > circle {
  stroke: #c2185b !important;
}

.pipeline-selection-highlight circle {
  stroke: none !important;
}

.pipeline-selection-highlight circle.white-highlight {
  fill: none !important;
  /* stroke: rgb(48, 48, 48) !important; */
  stroke: #c2185b !important;
  stroke-width: 4;
}

.svgResultStatus > circle.unknown {
  fill: #aaa !important;
  fill-opacity: 0.7;
}

.pipeline-line,
.pipeline-line path {
  stroke: #333333 !important;
}

marker path {
  fill: #333333 !important;
  stroke: #333333 !important;
}

.strike {
  text-decoration: line-through;
}

.spinning {
  animation: spin 8000ms linear infinite;
}

@keyframes spin {
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
    -webkit-transform-origin: center center;
    transform-origin: center center;
  }
}

.v-list-item:hover .v-btn {
  opacity: 1 !important;
}
</style>

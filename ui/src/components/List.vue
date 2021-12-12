<template>
  <v-row class="fill-height ma-0">
    <v-col cols="3" class="listnav" style="">
      <v-list nav color="background">
        <v-list-item
          v-if="showNew && canWrite"
          :to="{ name: routername, params: { id: 'new' } }"
          class="mt-4 mx-4 text-center newbutton"
        >
          <v-list-item-content>
            <v-list-item-title>
              <v-icon small>mdi-plus</v-icon> New {{ singular }}
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-subheader class="pl-4">{{ plural }}</v-subheader>
        <v-list-item
          v-for="item in (items ? items : [])"
          :key="item[itemid]"
          link
          :to="{ name: routername, params: { id: item[itemid] } }"
          class="mx-2"
        >
          <v-list-item-content>
            <v-list-item-title>
              {{ item[itemname] }}
            </v-list-item-title>
          </v-list-item-content>

          <v-list-item-action v-if="deletable && canWrite">
            <v-icon @click="askDelete(item[itemid])" class="fader">
              mdi-close
            </v-icon>
          </v-list-item-action>
        </v-list-item>
      </v-list>
    </v-col>
    <v-col cols="9">
      <router-view></router-view>
    </v-col>
    <v-dialog v-model="dialog" persistent max-width="400">
      <v-card>
        <v-card-title> Delete {{ singular }} {{ deleteName }} ? </v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="error" text @click="dialog = false">Cancel</v-btn>
          <v-btn color="success" outlined @click="deleteItem(deleteName)">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";

interface State {
  dialog: boolean;
  deleteName?: string;
}

export default Vue.extend({
  name: "List",
  components: {  },
  props: {
    items: {
      type: Array,
      required: true
    },
    routername: {
      type: String,
      required: true
    },
    itemid: {
      type: String,
      required: true
    },
    itemname: {
      type: String,
      required: true
    },
    singular: {
      type: String,
      required: true
    },
    plural: {
      type: String,
      required: true
    },
    showNew: {
      type: Boolean,
      default: true
    },
    deletable: {
      type: Boolean,
      default: true
    },
    writepermission: {
      type: String,
      required: true,
    }
  },
  data: (): State => ({
    dialog: false,
    deleteName: undefined,
  }),
  computed: {
    canWrite: function (): boolean {
      return this.hasRole(this.writepermission);
    },
  },
  methods: {
    askDelete(name: string) {
      this.deleteName = name;
      this.dialog = true;
    },
    deleteItem(deleteName: string) {
      this.$emit('delete', deleteName);
      this.dialog = false;
    },
    hasRole: function (s: string): boolean {
      if (this.$store.state.user.roles) {
        return this.lodash.includes(this.$store.state.user.roles, s);
      }
      return false;
    }
  }
});
</script>

<style>
.listnav {
  border-right: 1px solid #e0e0e0;
}
.theme--dark .listnav {
  border-right: 1px solid #393a3f;
}

.newbutton {
  background: #e0e0e0;
}
.theme--dark .newbutton {
  background: #424242;
}

.v-list-item .fader {
  opacity: 0 !important;
}

.v-list-item:hover .fader {
  opacity: 1 !important;
}
</style>

<template>
  <div>
    <div v-if="editoruserdata === undefined" class="text-sm-center py-16">
      <v-progress-circular
          indeterminate
          color="primary"
          :size="70"
          :width="7"
          class="align-center"
      >
      </v-progress-circular>
    </div>
    <div v-else class="pa-8">
      <v-form>
        <v-row>
          <v-col cols="4" class="d-flex flex-column align-center">
            <v-avatar v-if="editoruserdata.image" size="128" class="mt-1">
              <img :src="editoruserdata.image" alt="userdata avatar" />
            </v-avatar>

            <v-file-input
                v-model="file"
                type="file"
                class="pt-2 flex-grow-0"
                style="width: 100%"
                accept="image/png, image/jpeg"
                label="Select Image"
                @change="change"
                :clearable="false"
            >
              <template v-slot:append-outer>
                <v-btn
                    v-if="showCrop"
                    rounded
                    small
                    color="accent"
                    @click="validate"
                >
                  <v-icon>
                    mdi-check
                  </v-icon>
                  Set
                </v-btn>
                <v-btn
                    v-if="!!editoruserdata.image"
                    small
                    rounded
                    color="error"
                    @click="
                    file = null;
                    editoruserdata.image = '';
                    showCrop = false;
                  "
                    class="ml-2"
                >
                  <v-icon>
                    mdi-close
                  </v-icon>
                  Clear
                </v-btn>
              </template>
            </v-file-input>

            <vue-cropper
                v-if="showCrop"
                ref="cropper"
                v-bind="{ aspectRatio: 1, autoCrop: true }"
                :src="imgSrc"
                alt="Avatar"
            />
          </v-col>
          <v-col cols="8">
            <v-text-field
                prepend-icon="mdi-account"
                label="Name"
                v-model="editoruserdata.name"
            ></v-text-field>
            <v-text-field
                prepend-icon="mdi-email"
                label="Email"
                v-model="editoruserdata.email"
            ></v-text-field>

            <v-text-field
                prepend-icon="mdi-timetable"
                label="Timeformat"
                v-model="editoruserdata.timeformat"
            ></v-text-field>

            <v-btn
                color="success"
                outlined
                @click="saveUserData"
                class="mt-6"
            >
              <v-icon>mdi-content-save</v-icon>
              Save
            </v-btn>
          </v-col>
        </v-row>
      </v-form>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { UserData } from "@/client";
import VueCropper from "vue-cropperjs";
import "cropperjs/dist/cropper.css";

interface State {
  tab: number;
  editoruserdata?: UserData;
  file: File | null;
  imgSrc: string | ArrayBuffer | null;
  showCrop: boolean;
}

export default Vue.extend({
  name: "UserDataEditor",
  data: (): State => ({
    tab: 0,
    editoruserdata: undefined,
    file: null,
    imgSrc: null,
    showCrop: false
  }),
  props: ['userdata'],
  components: {
    VueCropper
  },
  methods: {
    saveUserData: function() {
      this.$emit("save", this.editoruserdata);
    },
    change: function() {
      if (!this.file) {
        this.imgSrc = null;
        return;
      }
      const reader = new FileReader();
      reader.onload = ticket => {
        if (ticket.target && this.$refs.cropper) {
          this.imgSrc = ticket.target.result;
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          let cropper: any = this.$refs.cropper;
          cropper.replace(this.imgSrc);
        }
      };
      reader.readAsDataURL(this.file);
      this.showCrop = true;
    },
    validate: function() {
      if (this.$refs.cropper && this.editoruserdata) {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        let cropper: any = this.$refs.cropper;
        this.editoruserdata.image = cropper
            .getCroppedCanvas({width: 128, height: 128})
            .toDataURL("image/png");
        // this.on.input(croppedImg)
        this.showCrop = false;
        // this.file = null
        // this.imgSrc = null
      }
    },
  },
  mounted() {
    this.editoruserdata = this.userdata;
  }
});
</script>

<style>
.theme--dark.v-tabs-items,
.theme--dark.v-tabs > .v-tabs-bar {
  background: none;
}
</style>

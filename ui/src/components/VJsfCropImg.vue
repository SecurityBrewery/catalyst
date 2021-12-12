<template>
  <v-input :value="value" class="vjsf-crop-img">
    <v-row class="mt-0 mx-0" align="center">
      <v-avatar v-if="value" size="128" class="mt-1">
        <img :src="value">
      </v-avatar>
      <v-file-input
          v-model="file"
          type="file"
          class="pt-2"
          accept="image/png, image/jpeg"
          placeholder="User Avatar"
          @change="change"
          :clearable="false"
      >
        <template v-slot:append-outer>
          <v-btn v-if="imgSrc" fab x-small color="accent" @click="validate">
            <v-icon>
              mdi-check
            </v-icon>
          </v-btn>
        </template>
      </v-file-input>
      <v-icon v-if="!!value" @click="on.input(null)">
        mdi-close
      </v-icon>
    </v-row>
    <vue-cropper
        v-if="file"
        ref="cropper"
        v-bind="cropperOptions"
        :src="imgSrc"
        alt="Avatar"
    />
  </v-input>
</template>

<script>
import VueCropper from 'vue-cropperjs'
import 'cropperjs/dist/cropper.css'

export default {
  components: { VueCropper },
  props: {
    value: { type: String, default: '' },
    on: { type: Object, required: true },
    cropperOptions: { type: Object, default: () => ({ aspectRatio: 1, autoCrop: true }) },
    size: { type: Number, default: 128 } // same as default v-avatar size
  },
  data: () => ({
    file: null,
    imgSrc: null
  }),
  computed: {},
  methods: {
    change() {
      if (!this.file) {
        this.imgSrc = null
        return
      }
      const reader = new FileReader()
      reader.onload = (ticket) => {
        this.imgSrc = ticket.target.result
        this.dialog = true
        this.$refs.cropper.replace(this.imgSrc)
      }
      reader.readAsDataURL(this.file)
    },
    async validate() {
      const croppedImg = this.$refs.cropper
          .getCroppedCanvas({ width: this.size, height: this.size })
          .toDataURL('image/png')
      this.on.input(croppedImg)
      this.file = null
      this.imgSrc = null
    }
  }
}
</script>

<style lang="css">
.vjsf-crop-img>.v-input__control>.v-input__slot {
  display: block;
}
</style>

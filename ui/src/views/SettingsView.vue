<script setup lang="ts">
import ColumnBody from '@/components/layout/ColumnBody.vue'
import ColumnBodyContainer from '@/components/layout/ColumnBodyContainer.vue'
import ColumnHeader from '@/components/layout/ColumnHeader.vue'
import TwoColumn from '@/components/layout/TwoColumn.vue'
import SettingsForm from '@/components/settings/SettingsForm.vue'
</script>

<template>
  <TwoColumn>
    <ColumnHeader title="Settings" show-sidebar-trigger />
    <ColumnBody>
      <ColumnBodyContainer small>
        <SettingsForm />
      </ColumnBodyContainer>
    </ColumnBody>
  </TwoColumn>
</template>

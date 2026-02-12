<style scoped>
.help {
  width: fit-content;
  background-color: rgba(var(--v-theme-background), 0.7);
  margin: auto;
}

.help div {
  margin: 5px 0px 5px 0px;
}

.help-icon {
  cursor: pointer;
  transition: all 0.3s ease-in-out;
}

.help-icon:hover {
  color: rgb(var(--v-theme-primary));
}
.help-tooltip :deep(.v-overlay__content) {
  border-radius: 25px;
  backdrop-filter: blur(5px);
  color: rgb(var(--v-theme-secondary));
  background-color: rgba(255, 255, 255, 0.1);
}

.cnd {
  font-weight: bold;
  border-radius: 10px;
  background-color: rgba(0, 0, 0, 0.1);
  padding: 3px 8px 3px 8px;
  border: 1px solid;
}

.exm {
  border: 1px solid rgba(83, 83, 83, 0.5);
  background-color: rgba(83, 83, 83, 0.3);
  border-radius: 7px;
  background-color: rgba();
  padding: 5px 10px 5px 10px;
}

.examples {
  line-height: 1.8rem;
}

.case {
  padding-bottom: 15px;
}
</style>

<template>
  <!-- Tooltip icon -->
  <v-tooltip class="help-tooltip" text="How to query" location="top">
    <template #activator="{ props }">
      <v-icon v-bind="props" class="help-icon" @click="handleOpenDialog">
        mdi-help-circle-outline
      </v-icon>
    </template>
  </v-tooltip>

  <!-- Help dialog -->
  <v-dialog v-model="dialog">
    <v-card class="help morphism d-flex flex-row pa-5" rounded="xl">
      <!-- Documentation -->
      <div>
        <h2>FEP conditions</h2>
        <h3>FEP +1</h3>
        <!-- FEP +1 -->
        <div>
          <span class="cnd str-t mr-2">str1</span><span class="cnd agi-t mr-2">agi1</span>
          <span class="cnd int-t mr-2">int1</span><span class="cnd con-t mr-2">con1</span>
          <span class="cnd prc-t mr-2">prc1</span><span class="cnd csm-t mr-2">csm1</span>
          <span class="cnd dex-t mr-2">dex1</span><span class="cnd wil-t mr-2">wil1</span>
          <span class="cnd psy-t">psy1</span> - FEP +1 value
        </div>

        <!-- FEP +2 -->
        <h3>FEP +2</h3>
        <div>
          <span class="cnd str-t mr-2">str2</span><span class="cnd agi-t mr-2">agi2</span>
          <span class="cnd int-t mr-2">int2</span><span class="cnd con-t mr-2">con2</span>
          <span class="cnd prc-t mr-2">prc2</span><span class="cnd csm-t mr-2">csm2</span>
          <span class="cnd dex-t mr-2">dex2</span><span class="cnd wil-t mr-2">wil2</span>
          <span class="cnd psy-t">psy2</span> - FEP +2 value
        </div>

        <div class="d-flex flex-row">
          <div>
            <h3>Available operators</h3>
            <div class="ga-1 mr-6">
              <span class="cnd mr-2">&gt;</span>
              <span class="cnd mr-2">&lt;</span>
              <span class="cnd mr-2">=</span>
              <span class="cnd mr-2">&gt;=</span>
              <span class="cnd mr-2">&lt;=</span>
            </div>
          </div>
          <div>
            <h3>Available types</h3>
            <div class="ga-1">
              <span class="cnd mr-2">number</span>
              <span class="cnd mr-2">percentage</span>
            </div>
          </div>
        </div>

        <v-divider class="mt-3 mb-1"></v-divider>

        <h2>Other conditions</h2>
        <!-- Other conditions -->
        <div>
          <span class="cnd">total</span> - FEP sum
          <div class="sub">
            └ Operators:
            <span class="cnd mr-2">&gt;</span>
            <span class="cnd mr-2">&lt;</span>
            <span class="cnd mr-2">=</span>
            <span class="cnd mr-2">&gt;=</span>
            <span class="cnd mr-2">&lt;=</span>
            Types:
            <span class="cnd mr-2">number</span>
            <span class="cnd mr-2">percentage</span>
          </div>
        </div>
        <div>
          <span class="cnd">hunger</span> - hunger % value
          <div>
            └ Operators:
            <span class="cnd mr-2">&gt;</span>
            <span class="cnd mr-2">&lt;</span>
            <span class="cnd mr-2">=</span>
            <span class="cnd mr-2">&gt;=</span>
            <span class="cnd mr-2">&lt;=</span>
            Types:
            <span class="cnd mr-2">number</span>
            <span class="cnd mr-2">percentage</span>
          </div>
        </div>
        <div>
          <span class="cnd">feph</span> - FEP / Hunger value
          <div>
            └ Operators:
            <span class="cnd mr-2">&gt;</span>
            <span class="cnd mr-2">&lt;</span>
            <span class="cnd mr-2">=</span>
            <span class="cnd mr-2">&gt;=</span>
            <span class="cnd mr-2">&lt;=</span>
            Types:
            <span class="cnd mr-2">number</span>
            <span class="cnd mr-2">percentage</span>
          </div>
        </div>
        <div>
          <span class="cnd">energy</span> - energy value
          <div>
            └ Operators:
            <span class="cnd mr-2">&gt;</span>
            <span class="cnd mr-2">&lt;</span>
            <span class="cnd mr-2">=</span>
            <span class="cnd mr-2">&gt;=</span>
            <span class="cnd mr-2">&lt;=</span>
            Types:
            <span class="cnd mr-2">number</span>
            <span class="cnd mr-2">percentage</span>
          </div>
        </div>
        <div>
          <span class="cnd">from</span> - made from ingredients
          <div>
            └ Operators:
            <span class="cnd mr-2">=</span>
            <span class="cnd mr-2">!=</span>
            Types:
            <span class="cnd mr-2">string</span>
            <span class="cnd mr-2">string list</span>
          </div>
        </div>
        <div>
          <span class="cnd">name</span> - item name
          <div>
            └ Operators:
            <span class="cnd mr-2">=</span><span class="cnd mr-2">!=</span>
            Types:
            <span class="cnd mr-2">string</span>
          </div>
        </div>
      </div>

      <v-divider class="mx-5" vertical></v-divider>

      <div class="d-flex flex-column" style="max-width: 500px">
        <div class="help">
          <h2 class="mb-4">Syntax</h2>
          <div class="case"><span class="exm">;</span> - conditions separator</div>
          <div class="case"><span class="exm">%</span> - percentage type</div>
          <div class="case">
            <span class="exm">10</span>, <span class="exm">3.5</span> - number type
          </div>
          <div class="case">
            <span class="exm">"..."</span>, <span class="exm">"...", "...", "..."</span> - string
            and string list type
          </div>
          <div class="case">
            <span class="exm">&gt;</span>, <span class="exm">&lt;</span>,
            <span class="exm">=</span>, <span class="exm">!=</span>, <span class="exm">&gt;=</span>,
            <span class="exm">&lt;=</span> - operators
          </div>
        </div>

        <div class="examples">
          <h2 class="mb-4">Examples</h2>
          <div class="case">
            <span class="exm">agi1>10;from=orca</span> - recipes from <b>"Orca meat"</b> with
            <b>"Agility +1"</b> more than <b>10 FEP</b>
          </div>
          <div class="case">
            <span class="exm">name=steak;str2>30%</span> - recipes with <b>"steak"</b> in name and
            <b>"Strength +2"</b> more than <b>30%</b>
          </div>
          <div class="case">
            <span class="exm">energy>=300;total&lt;10</span> - recipes with <b>energy</b> greater or
            equal to <b>300</b> and <b>total FEP</b> lower than <b>10</b>
          </div>
          <div class="case">
            <span class="exm">from!="dill","pepper","sage"</span> - recipes without <b>"Dill"</b>,
            <b>"Pepper"</b> and <b>"Sage"</b>
          </div>
        </div>
      </div>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const dialog = ref(false)

function handleOpenDialog() {
  dialog.value = true
}
</script>

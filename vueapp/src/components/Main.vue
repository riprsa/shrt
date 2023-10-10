<template>
  <div class="p-8 bg-neutral rounded-lg shadow-md max-w-md">

    <div v-if="!drawResult">
      <h1 class="text-2xl font-bold mb-4">Hello, user!</h1>
      <div class="w-full max-w-xs space-y-4">
        <input type="url" v-model="url" placeholder="Paste your link here" class="input input-primary w-full" required />
        <button @click="submit" class="btn btn-primary w-full">Shrt!</button>
      </div>
    </div>

    <div v-if="drawResult">
      <div v-if="short != '' && errorMessage == ''">
        <h1 class="text-2xl font-bold mb-4">URL Shortener Result</h1>
        <div class="mb-4">
          <span id="shortened-url" class="w-full p-2 rounded-md border border-gray-300 block">http://localhost/{{ short
          }}</span>
        </div>
        <div class="flex justify-between">
          <button @click="nullVars" class="btn btn-primary">Back</button>
          <button @click="copyURL" class="btn btn-primary">Copy link</button>
        </div>
      </div>
      <div v-if="short == '' && errorMessage != ''">
        <h1 class="text-2xl font-bold mb-4">URL Shortener Error</h1>
        <div class="mb-4">
          <label for="shortened-url" class="block text-gray-600 text-sm mb-2">Error:</label>
          <span id="shortened-url" class="w-full p-2 rounded-md border border-gray-300 block">{{ errorMessage }}</span>
        </div>
        <div class="flex justify-between">
          <button @click="nullVars" class=" btn btn-primary">Back</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

export default {
  data() {
    return {
      url: '',

      drawResult: false,

      errorMessage: '',

      short: ''
    }
  },
  methods: {
    nullVars() {
      this.drawResult = false;
      this.errorMessage = '';
      this.short = '';
      this.url = '';
    },

    async copyURL() {
      try {
        await navigator.clipboard.writeText("https://s.x16.me/" + this.short);
        // alert('Copied'); // TODO: may be alerts?
      } catch ($e) {
        // ignore
      }
    },

    async submit() {
      const rawResponse = await fetch('https://s.x16.me/api/short', {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ url: this.url })
      });

      if (rawResponse.status !== 200) {
        rawResponse.text().then((text) => {
          if (text === '') {
            text = 'Internal server error';
          }
          this.errorMessage = text;
          this.short = '';
          this.url = '';
        }).catch((e) => {
          console.error("strange error ->", e);
        });
      } else {
        rawResponse.json().then((data) => {
          console.log("data ->", data);
          this.short = data.short;

          this.errorMessage = '';
        }).catch((e) => {
          console.error("strange error ->", e);
        })
      }

      this.drawResult = true;
    }
  },

}
</script>
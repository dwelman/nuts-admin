<template>
  <div>

    <h1 class="mb-2">Discovery Services</h1>

    <ErrorMessage v-if="fetchError" :message="fetchError" :title="'Could not fetch Discovery Services'"/>

    <section>
      <p>Select a Discovery Service</p>
      <select v-on:change="viewService" id="discoveryServiceSelect">
        <option :value="service.id" v-for="service in services" :key="service.id">{{ service.id }}</option>
      </select>
    </section>

    <div v-if="selectedService">
      <section>
        <div>
          <label>Endpoint</label>
          <code>{{ selectedService.endpoint }}</code>
        </div>
      </section>
      <section>
        <header>Credential requirements</header>
        <div v-for="inputDescriptor in selectedService.presentation_definition.input_descriptors"
             :key="inputDescriptor.id" class="p-2 border-solid border-2 border-gray-400 rounded-md">
          <div>
            <label>Type:</label>
            <code>
              {{
                inputDescriptor.constraints.fields.find(f => f.path.includes('$.type') && f.filter && f.filter.type === "string").filter.const
              }}
            </code>
          </div>
        </div>
      </section>
      <section>
        <header>Search</header>
        <div class="search-quick-params">
          Add parameter:
          <a v-on:click="addSearchParam('id', 'did:web:example.com#abc')">Credential ID</a>
          <a v-on:click="addSearchParam('issuer', 'did:web:example.com')">Issuer DID</a>
          <a v-on:click="addSearchParam('credentialSubject.id', 'did:web:example.com#holder')">Holder DID</a>
          <a v-on:click="addSearchParam('', 'value')">Other</a>
        </div>
        <table class="w-full">
          <thead>
          <tr>
            <th class="w-1/3">Key</th>
            <th>Value</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="(param, idx) in searchParams" :key="'search-' + idx">
            <td><input type="text" v-model="param.key" v-on:keyup="search" v-on:blur="search" :id="'search-key-' + idx"></td>
            <td><input type="text" v-model="param.value" v-on:keyup="search" :placeholder="param.placeholder" v-on:blur="search" :id="'search-value-' + idx"></td>
          </tr>
          </tbody>
        </table>

        <div v-if="searchResults.length > 0">
          <label>Results</label>
          <table class="table w-full">
            <thead>
            <tr>
              <th class="thead">Subject DID</th>
              <th v-for="field in Object.keys(searchResults[0].fields)" class="thead">{{ field }}</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="result in searchResults" :key="'result-' + result.id">
              <td :key="result.id">{{ result.subject_id }}</td>
              <td v-for="field in Object.keys(result.fields)">{{ result.fields[field] }}</td>
            </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>
  </div>
</template>
<style>
.search-quick-params a {
  margin-left: 5px;
  cursor: pointer;
}
</style>
<script>
import ErrorMessage from "../components/ErrorMessage.vue";

export default {
  components: {ErrorMessage},
  data() {
    return {
      fetchError: '',
      services: [],
      selectedService: undefined,
      searchParams: [],
      searchResults: [],
    }
  },
  mounted() {
    this.$api.get('api/proxy/internal/discovery/v1')
        .then(data => {
          this.services = data
          // Select the first one by default
          if (this.services.length > 0) {
            this.viewService(this.services[0].id)
          }
        })
        .catch(response => {
          this.fetchError = response.statusText
          this.services = []
        })
  },
  methods: {
    viewService(id) {
      this.searchResults = []
      this.searchParams = []
      this.selectedService = this.services.find(s => s.id === id)
    },
    addSearchParam(key, placeholder) {
      this.searchParams.push({key: key, value: '', placeholder: placeholder})
      this.search()
    },
    search() {
      this.searchResults = []
      let entries = this.searchParams.filter(c => c.key && c.value).map(c => [c.key, c.value]);
      if (entries.length === 0) {
        return
      }
      const query = new URLSearchParams(entries);
      this.$api.get('api/proxy/internal/discovery/v1/' + this.selectedService.id + '?' + query.toString())
          .then(data => {
            // fields can be null, initialize with empty object with so. Otherwise, the UI crashes.
            data.forEach(r => r.fields = r.fields || {})
            this.searchResults = data
          })
          .catch(response => {
            this.fetchError = response.statusText
          })
    }
  }
}
</script>

<style>
section {
  margin-top: 20px;
}
</style>

---
sidebar_position: 1
---

# Getting the number from backend service

For frontend Mify uses [NuxtJS](https://nuxtjs.org/) template which is based on
Vue and provides a nice structure for the project. NuxtJS also provides great
support for Server Side + Client Side rendering. Here is a good description of page loading lifecycle: https://nuxtjs.org/docs/concepts/nuxt-lifecycle.

Let's edit the index page and add call to our backend (`js-services/counter-display/pages/index.vue`):

```html
<template>
  <div>
  <Sample/>
  Counter: {{counter}}
  </div>
</template>

<script>

export default {
  name: 'Index',
  data () {
    return {
      counter: 0,
    }
  },
  async fetch() {
    let countingBackend = this.$mifyContext.clients.countingBackend();
    try {
        var resp = await countingBackend.counterNextGet();
        this.counter = resp.number;
    } catch (e) {
        console.log(e);
    }
  }
}
</script>
```

We're added call to backend to the fetch stage and after receiving the number
we add it to the component data which can be accesses in a template.

Now we're ready for final testing.

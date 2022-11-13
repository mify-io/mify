---
sidebar_position: 2
---

# Building and Testing With Backend

To build frontend first we need to go to `js-services/counter-display` directory and install dependencies:
```
$ cd js-services/counter-display
$ yarn install
```

*Note: if you're using Ubuntu 20.04, this command won't work, because default
node is too old, please install node >= 18.12.1.*

Then we can run it:

```
$ yarn dev
```

*Note: Don't forget to keep backend running using `go run` command from [Building and Testing](/docs/guides/create-service/building-and-testing) backend section.*

Frontend service should be available at: [http://localhost:3000](http://localhost:3000])

You should see counter updating with every page refresh.

![](/img/docs/frontend-result.png)

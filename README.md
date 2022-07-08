k6-moonphase
============

An extension for [k6](https://k6.io) that provides current moon phase to your test scripts.

**This project has been created in educational purposes and you should not use it in production.**

Installation
------------

Use `xk6 build --with github.com/andrewslotin/k6-moonphase` to build a version of `k6` with this extension enabled.

Usage
-----

`k6-moonphase` uses [Stormglass API](https://stormglass.io/) to fetch the current moon phase in the given location. To access it, you need an API key that can be obtained from Stormglass upon registration. This key has to be provided to `k6` via the `STORMGLASS_API_KEY=` environment variable.

Here is an example of how to obtain the moon phase inside of your test script:

```javascript
import moonphase from 'k6/x/moonphase';

export default function () {
    const moonPhase = moonphase.current(52.52437, 13.41053);
    console.info(`Current moon phase is ${moonPhase.name}`)

    // ...
}
```

Here `52.52437, 13.41053` are the latitude and the longitude of the place.

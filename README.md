# jenkins-queue-health

This is currently under development, but the idea is to create a tool to
investigate the health of a build pipeline.  Failures happen, sometimes
because of the code checked in, and other times for external reasons.

The idea is to build a tool to determine the causes of failure so that
common or spurious failures can be dealt with to improve the realibility
of the pipeline.

## Building

On \*nix platforms with Make use the `Makefile`.

```
make && make test && sudo make install
```

Failing that you can read that file for the commands used to build and test.

## Running

```
jenkins-queue-health -url https://jenkins -project gerrit -user uname -password apitoken
```

Currently (at the time of writing this document) it's outputting json with
basic build info console logs for failing logs.

This is subject to change.  I'm still figuring out which way to go with this
tooling.

So far I've been using `jq` to analyse the output to look at problems, but the
intention is to move a large chunk of that into the tooling.

```
jq '.[] | select(.log | contains("Solr request failed - Timed out while waiting for socket to become ready for reading")) | { builtOn: .builtOn, timestamp: (.timestamp / 1000 | strftime("%Y-%m-%d %H:%M:%S")), build: .fullDisplayName }' "$1"
```

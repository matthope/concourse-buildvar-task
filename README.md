# concourse-buildvar-task

This image, intended on being used as a task in [Concourse](https://concourse-ci.org/), will output timestamps and environment variables (as a YAML file), suitable for use with `load_var`.

## Usage

```yaml
  plan:
    - task: buildvars
      image_resource:
        type: registry-image
        source: { repository: ghcr.io/matthope/concourse-buildvar-task }
      outputs:
        - name: buildvars
      ensure:
        load_var: buildvars
        file: buildvars/buildvars.yaml
        reveal: true

    - task: thing
      params:
        time: ((.:buildvars.time.rfc3339))
```

## Behind the Scenes

This task will iterate over all the directories in the current working directory, and create a file named "buildvars.yaml" in each.

This is so no configuration is required - the outputs can be named anything required.
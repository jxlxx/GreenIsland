# Green Island


## Getting Started


1. Have NATS running:

```
make up
```

2. Initialize the world:

(all this is doing is creating the kv buckets atm)

```
make init
```

3. Run the game:
```
make run-world
```


4. (optional) Connect to the natsbox

```
 make nats-box 
```


## Tips 

1. You can change the duration of a Game Hour in `world/world.go`. It is hardcoded in the `New() *World` function.



# Go pipe

`pipe` is a Go library used to make powerful data pipelines that are safe,
massively concurrent and reusable. Low level aspects of the implementation are taken cared of.

## How to

A `pipeline` is made of multiple `stages`, linked together through `channels`. Creating a `pipeline` is done as follows:
- Create a `pipeline` object.
- Create multiple `stage` objects.
- Add the `stages` to  the `pipeline`.
- Link the `stages` together through their `channels` via the `pipeline` object.
- Compile the `pipeline`, and start it.

## Example

To create a `stage`:
- Create a type that extends the `pipe.Stage` type.
- Overload the method `Pipeline(ctx context.Context)`.
- The `channels` used by the `Pipeline(ctx context.Context)` method must be store as fields of the extended type 
with type `pipe.Chan[T]`, where `T` is the type of the `channel`. These fields must be public.

For instance, here is a simple `stage` that multiplies its input by `2` before sending to its output:
```go
package main

type Multiplier struct {
	*pipe.Stage
	Input pipe.ReceiveChannel
	Output pipe.ReceiveChannel
}

func newMultiplier() *Multiplier {
	return &Printer{
		Stage: pipe.NewStage(),
		Input: pipe.NewChannel[int](), 
		Output: pipe.NewChannel[int](),
	}
}

func (m *Multiplier) Pipeline(ctx context.Context) {
	for {
		m.Output.Send(m.Input.Receive(ctx) * 2)
	}
}
```

As you can see, there is no need to handle the creation and destruction of channels, or to use `select` statements.

## Security

The library checks for issues automatically. Whenever a `stage` is added to a `pipeline`, it becomes the 
`pipeline`'s responsibility. The `pipeline` will make sure that the `stage` output channels are consumed by another `stage`,
and will refuse to compile if it is not. This is to prevent a `stage` from being blocked by trying to write to a `channel` 
that is full.

Note that this does not apply to the inputs of a `stage`: it is not considered as
an error if an `stage` input `channel` is not set to receive values. This is because a `stage` will never block if nothing is sent
to an input `channel`.

## Scalability

By default, the framework will create a single instance of each `stage` you define; this is the `pipe.NoScalability` flag.

Another flag is `pipe.MaxProcScalability`, which will instantiate `N` instances of each `stage` where `N` is equal to the
`GOMAXPROCS` environment variable.

Other scalability flags will be developed in the future. 

## Library

A set of reusable `stage` is provided as a library. Feel free to use them in your projects.
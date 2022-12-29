# Go pipe

`pipe` is a Go pipeline library used to make powerful data pipelines that are safe,
massively concurrent and reusable. Low level aspects of the implementation are taken
care of.

## How to

A `pipeline` is made of multiple `stage`, linked together through `channel`. Creating a `pipeline` is as follow:
- Create a `pipeline` object.
- Create multiple `stages` objects.
- Add `stages` to  the `pipeline`.
- Link `stages` together through their `channels` via the `pipeline` object.
- Compile the `pipeline`, and then start it.

## Example

To create a stage, do as follows:
- Create a type that composed with the `pipe.Stage` struct.
- Overload the method `Pipeline(ctx context.Context)`.
- The channels used by the `Pipeline(ctx context.Context)` method must be store as fields with type `pipe.Chan[T]`, where
`T` is the type of the channel. These fields must be public.

For instance, here is a simple stage that multiplies its input by `2` before sending to its output:
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
		m.Output.Send(ctx, m.Input.Receive(ctx) * 2)
	}
}
```

As you can see, there is no need to handle the creation and destruction of channels, or to use `select` statements.

## Scalability

By default, `pipe` will create a single instance of each stage you define; this is the `pipe.NoScalability` flag.

Another mode is `pipe.MaxProcScalability`, which will instantiate `N` instance of each stage where `N` is equal to the
`GOMAXPROCS` environment variable.

Other scalability mode will be later developed. 

## Library

A set of reusable `stage` is provided as a library. Feel free to use them or get inspiration from them.
package glumac

import (
	"github.com/teivah/gosiris/gosiris"
)

func helloActor() {
	gosiris.InitActorSystem(gosiris.SystemOptions{
		ActorSystemName: "HelloWorld",
	})

	pa := gosiris.Actor{}
	defer pa.Close()

	ca := gosiris.Actor{}
	defer ca.Close()

	ca.React("message", func(context gosiris.Context) {
		context.Self.LogInfo(context, "Received %v\n", context.Data)
	})

	gosiris.ActorSystem().RegisterActor("pa", &pa, nil)
	gosiris.ActorSystem().SpawnActor(&pa, "ca", &ca, nil)

	paRef, _ := gosiris.ActorSystem().ActorOf("pa")
	caRef, _ := gosiris.ActorSystem().ActorOf("ca")

	caRef.Tell(gosiris.EmptyContext, "message", "Hello, there, actor!", paRef)
}

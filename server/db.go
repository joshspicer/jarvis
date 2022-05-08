package main

func GetTrustedActors() []Actor {
	// TODO: example - actually pull from Cosmos
	testActor := Actor{
		id:        1,
		name:      "First",
		isTrusted: true,
		secret:    "AAAAAAAA-74C4-4039-A452-C5B4D1B59629_A25B82FA-C13E-4C2E-9659-23EFAF78D7DC",
	}

	testActor2 := Actor{
		id:        2,
		name:      "Second",
		isTrusted: true,
		secret:    "54887CAC-74C4-4039-A452-C5B4D1B59629_A25B82FA-C13E-4C2E-9659-23EFAF78D7DC",
	}

	return []Actor{
		testActor,
		testActor2,
	}
}

package merge

func convertoFiller(item Item, st State) filler {

	//
	return filler{
		item: item,
		// state: *state,
	}
}

// func convertoDataState(st State) state {

// }

func convertoDistcode(sc Sitecode) Distcode {
	scidx := sc.Index()
	// dc := (Distcode(sc.Layer()) << 56) | (Distcode(scidx/distRecordCount) << distRecordCount)
	dc := encodeBaseDistcode(sc.Layer(), uint8(scidx/distRecordCount))
	dc |= 1 << (scidx % distRecordCount)
	return dc
}

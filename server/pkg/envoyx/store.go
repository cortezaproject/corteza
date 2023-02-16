package envoyx

import "context"

func (svc *service) decodeStore(ctx context.Context, p DecodeParams) (nn NodeSet, err error) {
	var aux NodeSet
	for _, d := range svc.decoders[DecodeTypeStore] {
		aux, err = d.Decode(ctx, p)
		if err != nil {
			return
		}
		nn = append(nn, aux...)
	}

	return
}

func (svc *service) encodeStore(ctx context.Context, dg *depGraph, p EncodeParams) (err error) {
	// Prepping
	//
	// @note this is ok for now but if we add things like importing into
	//       multiple Cortezas at the same time, this won't be ok and each
	//       encoder should get it's own thing
	for rt, nn := range depNodesByResourceType(dg.allNodes()...) {
		for _, e := range svc.preparers[EncodeTypeStore] {
			err = e.Prepare(ctx, p, rt, OmitPlaceholderNodes(unpackDepNodes(nn...)...))
			if err != nil {
				return
			}
		}
	}

	// Encoding
	for rt, nn := range NodesByResourceType(dg.Roots()...) {
		for _, se := range svc.encoders[EncodeTypeStore] {
			err = se.Encode(ctx, p, rt, OmitPlaceholderNodes(nn...), dg)
			if err != nil {
				return
			}
		}
	}

	return
}

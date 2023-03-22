package envoyx

import (
	"context"
	"sort"
)

func (svc *Service) decodeStore(ctx context.Context, p DecodeParams) (nn NodeSet, err error) {
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

func (svc *Service) encodeStore(ctx context.Context, dg *DepGraph, p EncodeParams) (err error) {
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
	nodes := NodesByResourceType(dg.Roots()...)
	for _, rt := range svc.sortResourceTypes(nodes) {
		nn := nodes[rt]
		for _, se := range svc.encoders[EncodeTypeStore] {
			err = se.Encode(ctx, p, rt, OmitPlaceholderNodes(nn...), dg)
			if err != nil {
				return
			}
		}
	}

	return
}

func (svc *Service) sortResourceTypes(nn map[string]NodeSet) (out []string) {
	for rt := range nn {
		out = append(out, rt)
	}
	sort.Strings(out)

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return
}

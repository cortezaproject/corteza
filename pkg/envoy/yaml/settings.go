package yaml

type (
	settings map[string]interface{}
)

//func (wrap settings) MarshalEnvoy() ([]envoy.Node, error) {
//	nn := make([]envoy.Node, 0, 1+len(wrap.modules))
//	nn = append(nn, &envoy.SettingsNode{Ns: wrap.res})
//
//	if tmp, err := wrap.modules.MarshalEnvoy(); err != nil {
//		return nil, err
//	} else {
//		nn = append(nn, tmp...)
//	}
//
//	// @todo rbac
//
//	//if tmp, err := wrap.rules.MarshalEnvoy(); err != nil {
//	//	return nil, err
//	//} else {
//	//	nn = append(nn, tmp...)
//	//}
//
//	return nn, nil
//}

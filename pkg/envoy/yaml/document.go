package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	// Document defines the supported yaml structure
	Document struct {
		compose      *compose
		messaging    *messaging
		roles        roleSet
		users        userSet
		applications applicationSet
		settings     settingSet
		rbac         rbacRuleSet

		cfg *EncoderConfig
	}
)

func (doc *Document) UnmarshalYAML(n *yaml.Node) (err error) {
	if err = n.Decode(&doc.compose); err != nil {
		return
	}

	if err = n.Decode(&doc.messaging); err != nil {
		return
	}

	if doc.rbac, err = decodeRbac(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "roles":
			return v.Decode(&doc.roles)

		case "users":
			return v.Decode(&doc.users)

		case "applications":
			return v.Decode(&doc.applications)

		case "settings":
			return v.Decode(&doc.settings)

		}

		return nil
	})
}

func (doc *Document) MarshalYAML() (interface{}, error) {
	var err error
	dn, _ := makeMap()

	if doc.compose != nil {
		doc.compose.EncoderConfig = doc.cfg

		cn, err := encodeNode(doc.compose)
		if err != nil {
			return nil, err
		}

		dn, _ = inlineContent(dn, cn)
	}

	if doc.messaging != nil {
		doc.messaging.EncoderConfig = doc.cfg

		cn, err := encodeNode(doc.messaging)
		if err != nil {
			return nil, err
		}

		dn, _ = inlineContent(dn, cn)
	}

	if doc.roles != nil {
		doc.roles.ConfigureEncoder(doc.cfg)

		dn, err = encodeResource(dn, "roles", doc.roles, doc.cfg.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if doc.users != nil && len(doc.users) > 0 {
		doc.users.ConfigureEncoder(doc.cfg)

		dn, err = encodeResource(dn, "users", doc.users, doc.cfg.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if doc.applications != nil && len(doc.applications) > 0 {
		doc.applications.ConfigureEncoder(doc.cfg)

		// Applications don't support map representation
		dn, err = addMap(dn,
			"applications", doc.applications,
		)
		if err != nil {
			return nil, err
		}
	}

	if doc.settings != nil {
		dn, err = encodeResource(dn, "settings", doc.settings, doc.cfg.MappedOutput, "name")
		if err != nil {
			return nil, err
		}
	}

	if doc.rbac != nil && len(doc.rbac) > 0 {
		// rbac doesn't support map representation
		m, err := encodeNode(doc.rbac)
		if err != nil {
			return nil, err
		}

		dn, err = inlineContent(dn, m)
		if err != nil {
			return nil, err
		}
	}

	return dn, nil
}

func (doc *Document) Decode(ctx context.Context) ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, 100)

	mm := make([]envoy.Marshaller, 0, 20)
	if doc.compose != nil {
		mm = append(mm, doc.compose)
	}
	if doc.messaging != nil {
		mm = append(mm, doc.messaging)
	}
	if doc.roles != nil {
		mm = append(mm, doc.roles)
	}
	if doc.users != nil {
		mm = append(mm, doc.users)
	}
	if doc.applications != nil {
		mm = append(mm, doc.applications)
	}
	if doc.settings != nil {
		for _, s := range doc.settings {
			mm = append(mm, s)
		}
	}
	if doc.rbac != nil {
		mm = append(mm, doc.rbac)
	}

	for _, m := range mm {
		if tmp, err := m.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

// Document utilities...

// AddComposeNamespace adds a new composeNamespace to the document
func (doc *Document) AddComposeNamespace(n *composeNamespace) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		doc.compose.Namespaces = make(composeNamespaceSet, 0, 3)
	}

	doc.compose.Namespaces = append(doc.compose.Namespaces, n)
}

// AddComposeModule adds a new composeModule to the document
func (doc *Document) AddComposeModule(m *composeModule) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Modules == nil {
		doc.compose.Modules = make(composeModuleSet, 0)
	}

	doc.compose.Modules = append(doc.compose.Modules, m)
}

// AddComposeRecord adds a new composeRecord to the document
func (doc *Document) AddComposeRecord(r *composeRecord) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Records == nil {
		doc.compose.Records = make(composeRecordSet, 0)
	}

	doc.compose.Records = append(doc.compose.Records, r)
}

// AddComposePage adds a new composePage to the document
func (doc *Document) AddComposePage(p *composePage) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Pages == nil {
		doc.compose.Pages = make(composePageSet, 0)
	}

	doc.compose.Pages = append(doc.compose.Pages, p)
}

// AddComposeChart adds a new composeChart to the document
func (doc *Document) AddComposeChart(c *composeChart) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Charts == nil {
		doc.compose.Charts = make(composeChartSet, 0)
	}

	doc.compose.Charts = append(doc.compose.Charts, c)
}

// AddMessagingChannel adds a new messagingChannel to the document
func (doc *Document) AddMessagingChannel(c *messagingChannel) {
	if doc.messaging == nil {
		doc.messaging = &messaging{}
	}
	if doc.messaging.Channels == nil {
		doc.messaging.Channels = make(messagingChannelSet, 0)
	}

	doc.messaging.Channels = append(doc.messaging.Channels, c)
}

// AddRole adds a new role to the document
func (doc *Document) AddRole(r *role) {
	if doc.roles == nil {
		doc.roles = make(roleSet, 0, 20)
	}

	doc.roles = append(doc.roles, r)
}

// AddUser adds a new user to the document
func (doc *Document) AddUser(u *user) {
	if doc.users == nil {
		doc.users = make(userSet, 0, 20)
	}

	doc.users = append(doc.users, u)
}

// AddApplication adds a new application to the document
func (doc *Document) AddApplication(a *application) {
	if doc.applications == nil {
		doc.applications = make(applicationSet, 0, 20)
	}

	doc.applications = append(doc.applications, a)
}

// AddSetting adds a new setting to the document
func (doc *Document) AddSetting(s *setting) {
	if doc.settings == nil {
		doc.settings = make([]*setting, 0, 100)
	}

	doc.settings = append(doc.settings, s)
}

// AddRbacRule adds a new rbacRule to the document
func (doc *Document) AddRbacRule(r *rbacRule) {
	if doc.rbac == nil {
		doc.rbac = make(rbacRuleSet, 0, 100)
	}

	doc.rbac = append(doc.rbac, r)
}

// NestComposeModule adds a new composeModule to the document under a specified namespace
func (doc *Document) NestComposeModule(ns string, m *composeModule) error {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		return composeNamespaceErrNotFound(ns)
	}

	return doc.compose.Namespaces.AddComposeModule(ns, m)
}

// NestComposePage adds a new composePage to the document under a specified namespace
func (doc *Document) NestComposePage(ns string, p *composePage) error {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		return composeNamespaceErrNotFound(ns)
	}

	return doc.compose.Namespaces.AddComposePage(ns, p)
}

// NestComposePageChild adds a new composePage to the document under a specified parent page
func (doc *Document) NestComposePageChild(parent string, p *composePage) error {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		return composeNamespaceErrNotFound(parent)
	}

	return doc.compose.Pages.AddComposePage(parent, p)
}

// NestComposeChart adds a new composeChart to the document under a specified namespace
func (doc *Document) NestComposeChart(ns string, c *composeChart) error {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		return composeNamespaceErrNotFound(ns)
	}

	return doc.compose.Namespaces.AddComposeChart(ns, c)
}

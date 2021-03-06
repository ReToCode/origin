// This file was automatically generated by informer-gen

package v1

import (
	internalinterfaces "github.com/openshift/origin/pkg/user/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Users returns a UserInformer.
	Users() UserInformer
}

type version struct {
	internalinterfaces.SharedInformerFactory
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory) Interface {
	return &version{f}
}

// Users returns a UserInformer.
func (v *version) Users() UserInformer {
	return &userInformer{factory: v.SharedInformerFactory}
}

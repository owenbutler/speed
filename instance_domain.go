package speed

import "errors"

// InstanceDomain defines the interface for an instance domain
type InstanceDomain interface {
	ID() uint32                    // unique identifier for the instance domain
	Name() string                  // name of the instance domain
	Description() string           // description for the instance domain
	AddInstance(name string) error // adds an instance to the indom
	InstanceCount() int            // returns the number of instances in the indom
}

// PCPInstanceDomainBitLength is the maximum bit length of a PCP Instance Domain
//
// see: https://github.com/performancecopilot/pcp/blob/master/src/include/pcp/impl.h#L102-L121
const PCPInstanceDomainBitLength = 22

// PCPInstanceDomain wraps a PCP compatible instance domain
type PCPInstanceDomain struct {
	id                          uint32
	name                        string
	instances                   map[uint32]*Instance // the instances for this InstanceDomain stored as a map
	offset                      int
	instanceOffset              int
	shortHelpText, longHelpText string
}

// NewPCPInstanceDomain creates a new instance domain or returns an already created one for the passed name
// NOTE: this is different from parfait's idea of generating ids for InstanceDomains
// We simply generate a unique 32 bit hash for an instance domain name, and if it has not
// already been created, we create it, otherwise we return the already created version
func NewPCPInstanceDomain(name, shortDescription, longDescription string) *PCPInstanceDomain {
	return &PCPInstanceDomain{
		id:            getHash(name, PCPInstanceDomainBitLength),
		name:          name,
		instances:     make(map[uint32]*Instance),
		shortHelpText: shortDescription,
		longHelpText:  longDescription,
	}
}

// AddInstance adds a new instance to the current PCPInstanceDomain
func (indom *PCPInstanceDomain) AddInstance(name string) error {
	h := getHash(name, 0)

	_, present := indom.instances[h]
	if present {
		return errors.New("Instance with same name already created for the InstanceDomain")
	}

	ins := newInstance(h, name, indom)
	indom.instances[h] = ins

	return nil
}

// ID returns the id for PCPInstanceDomain
func (indom *PCPInstanceDomain) ID() uint32 { return indom.id }

// Name returns the name for PCPInstanceDomain
func (indom *PCPInstanceDomain) Name() string { return indom.name }

// InstanceCount returns the number of instances in the current instance domain
func (indom *PCPInstanceDomain) InstanceCount() int { return len(indom.instances) }

// Description returns the description for PCPInstanceDomain
func (indom *PCPInstanceDomain) Description() string {
	s, l := indom.shortHelpText, indom.longHelpText
	if l != "" {
		return s + "\n\n" + l
	}
	return s
}

func (indom *PCPInstanceDomain) Offset() int { return indom.offset }

func (indom *PCPInstanceDomain) SetOffset(offset int) { indom.offset = offset }

func (indom *PCPInstanceDomain) String() string {
	s := "InstanceDomain: " + indom.name
	if len(indom.instances) > 0 {
		s += "["
		for _, i := range indom.instances {
			s += i.name + ","
		}
		s += "]"
	}
	return s
}

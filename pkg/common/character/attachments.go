package character

type AttachmentType int

const (
	AttachmentShadow = AttachmentType(iota)
	AttachmentBody
	AttachmentHead
	AttachmentTop
	AttachmentMid
	AttachmentBottom

	NumAttachments
)

func (e AttachmentType) String() (att string) {
	switch e {
	case AttachmentShadow:
		att = "AttachmentShadow"
	case AttachmentBody:
		att = "AttachmentBody"
	case AttachmentHead:
		att = "AttachmentHead"
	case AttachmentTop:
		att = "AttachmentTop"
	case AttachmentMid:
		att = "AttachmentMid"
	case AttachmentBottom:
		att = "AttachmentBottom"
	default:
		panic("unsupported attachment type")
	}

	return att
}

func Attachments() []AttachmentType {
	return []AttachmentType{
		AttachmentShadow,
		AttachmentBody,
		AttachmentHead,
		AttachmentTop,
		AttachmentMid,
		AttachmentBottom,
	}
}

package actionindex

type Type int

const (
	Idle            Type = 0
	Walking         Type = 8
	Sitting         Type = 16
	PickingItem     Type = 24
	StandBy         Type = 32
	Attacking1      Type = 40
	ReceivingDamage Type = 48
	Freeze1         Type = 56
	Dead            Type = 65
	Freeze2         Type = 72
	Attacking2      Type = 80
	Attacking3      Type = 88
	CastingSpell    Type = 96
)

package internal

import "github.com/rprtr258/sixty"

// TODO: funcs to encode some primitive/logical types

// TODO: name fields as in comments
var (
	BttvEmotesDBSchemaString = `[](std:string,[]std:string)`
	BttvEmotesDBSchema       = sixty.List{Schema: sixty.Product{
		sixty.String{}, // "emote"
		sixty.List{ // "ids"
			// TODO: change to const(?) length id
			Schema: sixty.String{},
		},
	}}
	BttvEmotesDBData = sixty.ListValue{sixty.ProductValue{
		sixty.StringValue("Kappa"),
		sixty.ListValue{
			sixty.StringValue("KappaOwner1"),
			sixty.StringValue("KappaOwner2"),
		},
	}, sixty.ProductValue{
		sixty.StringValue("PogChamp"),
		sixty.ListValue{
			sixty.StringValue("PogChampOwner1"),
			sixty.StringValue("PogChampOwner2"),
		},
	}}
	BttvEmotesDBBytes = []byte{
		0x2,
		0x5, 0x4b, 0x61, 0x70, 0x70, 0x61, // "Kappa"
		0x2,
		0xb, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "KappaOwner1"
		0xb, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "KappaOwner2"
		0x8, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, // "PogChamp"
		0x2,
		0xe, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "PogChampOwner1"
		0xe, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "PogChampOwner2"
	}
	BttvEmotesDBJson = `[{"0":"Kappa","1":["KappaOwner1","KappaOwner2"]},{"0":"PogChamp","1":["PogChampOwner1","PogChampOwner2"]}]`
)

var (
	VKInfoSchema = sixty.Coproduct{
		sixty.Product{ // "group"
			sixty.Uint{},   // "id"
			sixty.String{}, // "alias"
			sixty.String{}, // "name"
		},
		sixty.Product{ // "user"
			sixty.Uint{},   //"id"
			sixty.String{}, // "first_name"
			sixty.String{}, // "last_name"
		},
	}
	VKInfoSchemaString = `[(uint,std:string,std:string),(uint,std:string,std:string)]`
)

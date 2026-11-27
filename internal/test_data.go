package internal

import (
	_ "embed"

	"github.com/rprtr258/sixty"
)

// TODO: funcs to encode some primitive/logical types

// TODO: name fields as in comments.
var (
	BttvEmotesDBSchema = sixty.List{Schema: sixty.Product{
		{"emote", sixty.String{}},
		{"ids", sixty.List{
			// TODO: change to const(?) length id
			Schema: sixty.String{},
		}},
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
)

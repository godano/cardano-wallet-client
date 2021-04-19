package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const LIST_WALLETS_SUCCESS = `[
	{
	  "id": "2512a00e9653fe49a44a5886202e24d77eeb998f",
	  "address_pool_gap": 20,
	  "balance": {
		"available": {
		  "quantity": 42000000,
		  "unit": "lovelace"
		},
		"reward": {
		  "quantity": 42000000,
		  "unit": "lovelace"
		},
		"total": {
		  "quantity": 42000000,
		  "unit": "lovelace"
		}
	  },
	  "assets": {
		"available": [
		  {
			"policy_id": "65ab82542b0ca20391caaf66a4d4d7897d281f9c136cd3513136945b",
			"asset_name": "",
			"quantity": 0
		  }
		],
		"total": [
		  {
			"policy_id": "65ab82542b0ca20391caaf66a4d4d7897d281f9c136cd3513136945b",
			"asset_name": "",
			"quantity": 0
		  }
		]
	  },
	  "delegation": {
		"active": {
		  "status": "delegating",
		  "target": "1423856bc91c49e928f6f30f4e8d665d53eb4ab6028bd0ac971809d514c92db1"
		},
		"next": [
		  {
			"status": "not_delegating",
			"changes_at": {
			  "epoch_number": 14,
			  "epoch_start_time": "2020-01-22T10:06:39.037Z"
			}
		  }
		]
	  },
	  "name": "Alan's Wallet",
	  "passphrase": {
		"last_updated_at": "2019-02-27T14:46:45.000Z"
	  },
	  "state": {
		"status": "ready"
	  },
	  "tip": {
		"absolute_slot_number": 8086,
		"slot_number": 1337,
		"epoch_number": 14,
		"time": "2019-02-27T14:46:45.000Z",
		"height": {
		  "quantity": 1337,
		  "unit": "block"
		}
	  }
	}
  ]`

func TestListWallets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, LIST_WALLETS_SUCCESS)
	}))
	defer ts.Close()

	client := NewClient()
	client.BaseURL = ts.URL

	results, err := client.ListWallets(context.Background())

	assert.Equal(t, len(results), 1, "Should only have one wallet")
	assert.Nil(t, err)

}

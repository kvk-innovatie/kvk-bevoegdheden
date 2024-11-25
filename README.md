KVK Extract (Golang lib)
------

GO library that fetches an extract from the KVK Dataservice and extracts all fields from it, that are relevant for a 'Bevoegdheid' or 'Machtiging'.

To call the getBevoegdheid function
```
import (
	kvkBevoegdheden "github.com/kvk-innovatie/kvk-bevoegdheden"
	"github.com/kvk-innovatie/kvk-bevoegdheden/models"
)

bevoegdheidResponse, err, rawResponse := kvkBevoegdheden.GetBevoegdheid(kvkNummer, *session.IdentityNP, os.Getenv("LOOKUP_CLIENTID"), os.Getenv("LOOKUP_CLIENTSECRET"), os.Getenv("LOOKUP_AUTHSERVER_URL"), enableCaching, env)
```
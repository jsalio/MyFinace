module Financial/persistence

go 1.23.9

require (
	Financial/Core v0.0.0
	github.com/supabase-community/postgrest-go v0.0.11
	github.com/supabase-community/supabase-go v0.0.4
)

replace Financial/Core => ../Core

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/supabase-community/functions-go v0.0.0-20220927045802-22373e6cb51d // indirect
	github.com/supabase-community/gotrue-go v1.2.0 // indirect
	github.com/supabase-community/storage-go v0.7.0 // indirect
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
)

package utils


const (
	
	//keys for watchers
	MAILGUN = "mailgun"
	SLACK   = "slack"
	INFOBIP = "infobip"
	SCYLLA  = "scylla"
	META    = "meta"
	WHMCS   = "whmcs"

	//config keys by watchers
	TOKEN          = "token"
	CHANNEL        = "channel"
	USERNAME       = "username"
	PASSWORD       = "password"
	APPLICATION_ID = "application_id"
	MESSAGE_ID     = "message_id"
	API_KEY        = "api_key"
	DOMAIN         = "domain"
	PIGGYBANKS     = "piggybanks"
	
	HOME           = "home"
	DIR            = "dir"
	SCYLLAHOST     = "scylla_host"
	SCYLLAKEYSPACE = "scylla_keyspace"
	
	EVENT_TYPE      = "event_type"
	ACCOUNT_ID     = "account_id"
	ASSEMBLY_ID    = "assembly_id"
	DATA		   = "data"
	CREATED_AT     = "created_at"

	//args for notification
	NILAVU    = "nilavu"
	LOGO      = "logo"
	NAME      = "name"
	VERTNAME  = "appname"
	TEAM      = "team"
	VERTTYPE  = "type"
	EMAIL     = "email"
	DAYS      = "days"
	COST      = "cost"
	STARTTIME = "starttime"
	ENDTIME   = "endtime"
	//STATUS    = "status"
	//DESCRIPTION = "description"
	
	EventMachine                 = "machine"
	EventContainer               = "container"
	EventBill                    = "bill"
	EventUser                    = "user"
	EventStatus                  = "status"
	
	BILLMGR = "bill"
	ADDONS  = "addons"
	
	PROVIDER        = "provider"
	PROVIDER_ONE    = "one"
	PROVIDER_DOCKER = "docker"

	LAUNCHING    = "launching"
	LAUNCHED     = "launched"
	BOOTSTRAPPED = "bootstrapped"
	BOOTSTRAPPING = "bootstrapping"
	STATEUP      = "stateup"
	RUNNING      = "running"
	STARTING     = "starting"
	STARTED      = "started"
	STOPPING     = "stopping"
	STOPPED      = "stopped"
	RESTARTING     = "restarting"
	RESTARTED      = "restarted"
	UPGRADED     = "upgraded"
	DESTROYING   = "destroying"
	NUKED        = "nuked"
	ERROR        = "error"
	
	COOKBOOKSUCCESS = "cookbook_success"
	COOKBOOKFAILURE = "cookbook_failure"
	AUTHKEYSSUCCESS = "authkeys_success"
	AUTHKEYSFAILURE = "authkeys_failure"
	INSTANCEIPSSUCCESS = "ips_success"
	INSTANCEIPSFAILURE = "ips_failure"
	
	CONTAINERNETWORKSUCCESS = "container_network_success"
	CONTAINERNETWORKFAILURE = "container_network_failure"
	
	// StatusLaunching is the initial status of a box
	// it should transition shortly to a more specific status
	StatusLaunching = Status(LAUNCHING)

	// StatusLaunched is the status for box after launched in cloud.
	StatusLaunched = Status(LAUNCHED)

	// StatusBootstrapped is the status for box after being booted by the agent in cloud
	StatusBootstrapped = Status(BOOTSTRAPPED)
	StatusBootstrapping = Status(BOOTSTRAPPING)

	// Stateup is the status of the which is moving up in the state in cloud.
	// Sent by vertice to gulpd when it received StatusBootstrapped.
	StatusStateup = Status(STATEUP)

	//fully up instance
	StatusRunning = Status(RUNNING)

	StatusStarting = Status(STARTING)
	StatusStarted  = Status(STARTED)

	StatusStopping = Status(STOPPING)
	StatusStopped  = Status(STOPPED)
	
	StatusRestarting = Status(RESTARTING)
	StatusRestarted  = Status(RESTARTED)

	StatusUpgraded = Status(UPGRADED)

	StatusDestroying = Status(DESTROYING)
	StatusNuked      = Status(NUKED)

	// StatusError is the status for units that failed to start, because of
	// a box error.
	StatusError = Status(ERROR)
	
	StatusCookbookSuccess = Status(COOKBOOKSUCCESS)
	StatusCookbookFailure = Status(COOKBOOKFAILURE)
	StatusAuthkeysSuccess = Status(AUTHKEYSSUCCESS)
	StatusAuthkeysFailure = Status(AUTHKEYSFAILURE)
	StatusIpsSuccess = Status(INSTANCEIPSSUCCESS)
	StatusIpsFailure = Status(INSTANCEIPSFAILURE)
	
	StatusContainerNetworkSuccess = Status(CONTAINERNETWORKSUCCESS)
	StatusContainerNetworkFailure = Status(CONTAINERNETWORKFAILURE)

	ONEINSTANCELAUNCHINGTYPE     = "compute.instance.launching"
	ONEINSTANCEBOOTSTRAPPINGTYPE = "compute.instance.bootstrapping"
	ONEINSTANCEBOOTSTRAPPEDTYPE  = "compute.instance.bootstrapped"
	ONEINSTANCESTATEUPTYPE       = "compute.instance.stateup"
	ONEINSTANCERUNNINGTYPE       = "compute.instance.running"
	ONEINSTANCELAUNCHEDTYPE      = "compute.instance.launched"
	ONEINSTANCEEXISTSTYPE        = "compute.instance.exists"
	ONEINSTANCEDESTROYINGTYPE    = "compute.instance.destroying"
	ONEINSTANCEDELETEDTYPE       = "compute.instance.deleted"
	ONEINSTANCESTARTINGTYPE      = "compute.instance.starting"
	ONEINSTANCESTARTEDTYPE       = "compute.instance.started"
	ONEINSTANCESTOPPINGTYPE      = "compute.instance.stopping"
	ONEINSTANCESTOPPEDTYPE       = "compute.instance.stopped"
	ONEINSTANCERESTARTINGTYPE    = "compute.instance.restarting"
	ONEINSTANCERESTARTEDTYPE     = "compute.instance.restarted"
	ONEINSTANCEUPGRADEDTYPE      = "compute.instance.upgraded"
	ONEINSTANCERESIZINGTYPE      = "compute.instance.resizing"
	ONEINSTANCERESIZEDTYPE       = "compute.instance.resized"
	ONEINSTANCEERRORTYPE         = "compute.instance.error"
	
	COOKBOOKSUCCESSTYPE          = "compute.instance.cookbook.download.success"
	COOKBOOKFAILURETYPE          = "compute.instance.cookbook.download.failure"
	AUTHKEYSSUCCESSTYPE          = "compute.instance.authkeys.success"
	AUTHKEYSFAILURETYPE          = "compute.instance.authkeys.failure"
	INSTANCEIPSSUCCESSTYPE       = "net.instance.ip.update.success"
	INSTANCEIPSFAILURETYPE       = "net.instance.ip.update.failure"
	
	CONTAINERNETWORKSUCCESSTYPE      = "net.container.ip.allocate.success"
	CONTAINERNETWORKFAILURETYPE      = "net.container.ip.allocate.failure"
	
)

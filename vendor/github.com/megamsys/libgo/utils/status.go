/*
** Copyright [2013-2016] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package utils

// Status represents the status of a unit in vertice
type Status string

func (s Status) String() string {
	return string(s)
}

func (s Status) Event_type() string {
	switch s.String() {
	case LAUNCHING:
		return ONEINSTANCELAUNCHINGTYPE
	case LAUNCHED:
		return ONEINSTANCELAUNCHEDTYPE
	case BOOTSTRAPPED:
		return ONEINSTANCEBOOTSTRAPPEDTYPE
	case STATEUP:
		return ONEINSTANCESTATEUPTYPE
	case RUNNING:
		return ONEINSTANCERUNNINGTYPE
	case STARTING:
		return ONEINSTANCESTARTINGTYPE
	case STARTED:
		return ONEINSTANCESTARTEDTYPE
	case STOPPING:
		return ONEINSTANCESTOPPINGTYPE
	case STOPPED:
		return ONEINSTANCESTOPPEDTYPE
	case UPGRADED:
		return ONEINSTANCEUPGRADEDTYPE
	case DESTROYING:
		return ONEINSTANCEDESTROYINGTYPE
	case NUKED:
		return ONEINSTANCEDELETEDTYPE
	case COOKBOOKSUCCESS: 
	    return 	COOKBOOKSUCCESSTYPE	
	case COOKBOOKFAILURE: 
	    return 	COOKBOOKFAILURETYPE	  
	case AUTHKEYSSUCCESS: 
	    return 	AUTHKEYSSUCCESSTYPE	
	case AUTHKEYSFAILURE: 
	    return 	AUTHKEYSFAILURETYPE	   
	case INSTANCEIPSSUCCESS: 
	    return 	INSTANCEIPSSUCCESSTYPE	
	case INSTANCEIPSFAILURE: 
	    return 	INSTANCEIPSFAILURETYPE	
	case CONTAINERNETWORKSUCCESS: 
	    return 	CONTAINERNETWORKSUCCESSTYPE	
	case CONTAINERNETWORKFAILURE: 
	    return 	CONTAINERNETWORKFAILURETYPE	               
	case ERROR:
		return ONEINSTANCEERRORTYPE
	default:
		return "arrgh"
	}
}

func (s Status) Description(name string) string {
	error_common := "Oops something went wrong on"
	switch s.String() {
	case LAUNCHING:
		return "your " + name + " machine is initializing.."
	case LAUNCHED:
		return "Machine " + name + " was initialized on cloud.."
	case BOOTSTRAPPED:
		return name + " was booted.."
	case STATEUP:
		return "moving up in the state of " + name + ".."
	case RUNNING:
		return name + " is running.."
	case STARTING:
		return "starting process initializing on " + name + ".."
	case STARTED:
		return name + " was started.."
	case STOPPING:
		return "stopping process initializing on " + name + ".."
	case STOPPED:
		return name + " was stopped.."
	case UPGRADED:
		return name + " was upgraded.."
	case DESTROYING:
		return "destroying process initializing on " + name + ".."
	case NUKED:
		return name + " was removed.."
	case COOKBOOKSUCCESS:
	    return "chef cookbooks are successfully downloaded.."	
	case COOKBOOKFAILURE:
	    return error_common + "downloading cookbooks on " + name + ".."	    
	case AUTHKEYSSUCCESS:
	    return "SSH keys are updated on your " + name	
	case AUTHKEYSFAILURE:
	    return error_common + "updating Ssh keys on " + name + ".."	  
	case INSTANCEIPSSUCCESS:
	    return "Private and public ips are updated on your " + name	
	case INSTANCEIPSFAILURE:
	    return error_common + "updating private and public ips on " + name + ".."	 
	case CONTAINERNETWORKSUCCESS:
	    return "Private and public ips are updated on your " + name	
	case CONTAINERNETWORKFAILURE:
	    return error_common + "updating private and public ips on " + name + ".."	               	
	case ERROR:
		return "oops something went wrong on " + name + ".."
	default:
		return "arrgh"
	}
}




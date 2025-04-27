export namespace model {
	
	export class Secret {
	    ID: number;
	    AccountName: string;
	    ServerName: string;
	    Code: string;
	
	    static createFrom(source: any = {}) {
	        return new Secret(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.AccountName = source["AccountName"];
	        this.ServerName = source["ServerName"];
	        this.Code = source["Code"];
	    }
	}

}


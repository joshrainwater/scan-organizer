export namespace scanorganizer {
	
	export class PreviewData {
	    preview: string;
	    previousRenamed: string[];
	    folders: string[];
	
	    static createFrom(source: any = {}) {
	        return new PreviewData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.preview = source["preview"];
	        this.previousRenamed = source["previousRenamed"];
	        this.folders = source["folders"];
	    }
	}

}


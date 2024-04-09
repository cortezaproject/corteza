//shim toallow node and native browser module path support with the same code
import Is from '../../@achrinza/strong-type/index.js';

const is=new Is;

class EventPubSub {
    constructor() {
        
    }

    on(type, handler, once=false) {
        is.string(type);
        is.function(handler);
        is.boolean(once);
        
        if(type=='*'){
            type=this.#all;
        }

        if (!this.#events[type]) {
            this.#events[type] = [];
        }

        handler[this.#once] = once;

        this.#events[type].push(handler);
        
        return this;
    }

    once(type, handler) {
        //sugar for this.on with once set to true 
        //so let that do the validation
        return this.on(type,handler,true);
    }

    off(type='*', handler='*') {
        is.string(type);
        
        if(type==this.#all.toString()||type=='*'){
            type=this.#all;
        }
        
        if (!this.#events[type]) {
            return this;
        }

        if (handler=='*') {
            delete this.#events[type];
            return this;
        }

        //If we are not removing all the handlers,
        //we need to know which one we are removing.
        is.function(handler);

        const handlers = this.#events[type];

        while (handlers.includes(handler)) {
            handlers.splice(
                handlers.indexOf( handler ),
                1
            );
        }

        if (handlers.length < 1) {
            delete this.#events[type];
        }

        return this;
    }

    emit(type, ...args) {
        is.string(type);
        
        const globalHandlers=this.#events[this.#all]||[];
        
        this.#handleOnce(this.#all.toString(), globalHandlers, type, ...args);
        
        if (!this.#events[type]) {
            return this;
        }

        const handlers = this.#events[type];        

        this.#handleOnce(type, handlers, ...args);

        return this;
    }

    reset(){
        this.off(this.#all.toString());
        for(let type in this.#events){
            this.off(type);
        }

        return this
    }

    get list(){
        return Object.assign({},this.#events);
    }

    #handleOnce=(type, handlers, ...args)=>{
        is.string(type);
        is.array(handlers);
        
        const deleteOnceHandled=[];

        for (let handler of handlers) {
            handler(...args);
            if(handler[this.#once]){
                deleteOnceHandled.push(handler);
            }
        }

        for(let handler of deleteOnceHandled){
          this.off(type,handler);
        }
    }

    #all =Symbol.for('event-pubsub-all')
    #once=Symbol.for('event-pubsub-once')

    #events={}
}

export {EventPubSub as default, EventPubSub};
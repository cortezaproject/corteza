class Fake{
    //fake class as fallback
}

class FakeCore{
    //fake class as fallback
}

class Is{
    constructor(strict=true){
        this.strict=strict;
    }

    //core
    throw(valueType,expectedType){
        let err=new TypeError;
        err.message=`expected type of ${valueType} to be ${expectedType}`;
        if(!this.strict){
            return false;
        }
        throw err;
    }

    typeCheck(value,type){
        if(typeof value === type){
            return true;
        }
        return this.throw(typeof value,type);
    }

    instanceCheck(value=new Fake, constructor=FakeCore){
        //console.log(value,constructor);
        if(value instanceof constructor){
            return true;
        }
        return this.throw(typeof value,constructor.name);
    }

    symbolStringCheck(value,type){
        if(Object.prototype.toString.call(value) == `[object ${type}]`){
            return true;
        }
        return this.throw(Object.prototype.toString.call(value),`[object ${type}]`);
    }

    compare(value,targetValue,typeName){
        if(value==targetValue){
            return true;
        }
        return this.throw(typeof value, typeName);
    }

    defined(value){
        const weakIs=new Is(false);
        if(weakIs.undefined(value)){
            return this.throw('undefined','defined');
        }

        return true;
    }

    any(value){
        return this.defined(value);
    }

    exists(value){
        return this.defined(value);
    }

    union(value,typesString){
        const types=typesString.split('|');
        const weakIs=new Is(false);
        let pass=false;
        let type='undefined';
        for(type of types){
            try{
                if(weakIs[type](value)){
                    pass=true;
                    break;
                }
            }catch(err){
                return this.throw(type,'a method available on strong-type');
            }
        }

        if(pass){
           return this[type](value);
        }

        return this.throw(typeof value, types.join('|'));

    }

    //unique checks
    finite(value){
        if(isFinite(value)){
            return true;
        }
        return this.throw(typeof value, 'finite');
    }

    NaN(value){
        if(!this.number(value)){
            return this.number(value);
        }

        if(isNaN(value)){
            return true;
        }

        return this.throw(typeof value, 'NaN');
    }

    null(value){
        return this.compare(value,null,'null');
    }

    //common sugar
    array(value){
        return this.instanceCheck(value,Array);
    }

    boolean(value){
        return this.typeCheck(value,'boolean');
    }

    bigInt(value){
        return this.typeCheck(value,'bigint');
    }

    date(value){
        return this.instanceCheck(value,Date);
    }

    generator(value){
        return this.symbolStringCheck(value,'Generator');
    }

    asyncGenerator(value){
        return this.symbolStringCheck(value,'AsyncGenerator');
    }

    globalThis(value){
        return this.compare(value,globalThis,'explicitly globalThis, not window, global nor self');
    }

    infinity(value){
        return this.compare(value,Infinity,'Infinity');
    }

    map(value){
        return this.instanceCheck(value,Map);
    }

    weakMap(value){
        return this.instanceCheck(value,WeakMap);
    }

    number(value){
        return this.typeCheck(value,'number');
    }

    object(value){
        return this.typeCheck(value,'object');
    }

    promise(value){
        return this.instanceCheck(value,Promise);
    }

    regExp(value){
        return this.instanceCheck(value,RegExp);
    }
    
    undefined(value){
        return this.typeCheck(value,'undefined');
    }

    set(value){
        return this.instanceCheck(value,Set);
    }

    weakSet(value){
        return this.instanceCheck(value,WeakSet);
    }
    
    string(value){
        return this.typeCheck(value,'string');
    }

    symbol(value){
        return this.typeCheck(value,'symbol');
    }

    //functions
    function(value){
        return this.typeCheck(value,'function');
    }

    asyncFunction(value){
        return this.symbolStringCheck(value,'AsyncFunction');
    }

    generatorFunction(value){
        return this.symbolStringCheck(value,'GeneratorFunction');
    }

    asyncGeneratorFunction(value){
        return this.symbolStringCheck(value,'AsyncGeneratorFunction');
    }

    //error sugar
    error(value){
        return this.instanceCheck(value,Error);
    }

    evalError(value){
        return this.instanceCheck(value,EvalError);
    }

    rangeError(value){
        return this.instanceCheck(value,RangeError);
    }

    referenceError(value){
        return this.instanceCheck(value,ReferenceError);
    }

    syntaxError(value){
        return this.instanceCheck(value,SyntaxError);
    }

    typeError(value){
        return this.instanceCheck(value,TypeError);
    }

    URIError(value){
        return this.instanceCheck(value,URIError);
    }    

    //typed array sugar
    bigInt64Array(value){
        return this.instanceCheck(value,BigInt64Array);
    }

    bigUint64Array(value){
        return this.instanceCheck(value,BigUint64Array);
    }

    float32Array(value){
        return this.instanceCheck(value,Float32Array);
    }

    float64Array(value){
        return this.instanceCheck(value,Float64Array);
    }

    int8Array(value){
        return this.instanceCheck(value,Int8Array);
    }

    int16Array(value){
        return this.instanceCheck(value,Int16Array);
    }

    int32Array(value){
        return this.instanceCheck(value,Int32Array);
    }

    uint8Array(value){
        return this.instanceCheck(value,Uint8Array);
    }

    uint8ClampedArray(value){
        return this.instanceCheck(value,Uint8ClampedArray);
    }
    
    uint16Array(value){
        return this.instanceCheck(value,Uint16Array);
    }

    uint32Array(value){
        return this.instanceCheck(value,Uint32Array);
    }

    //buffers
    arrayBuffer(value){
        return this.instanceCheck(value,ArrayBuffer);
    }

    dataView(value){
        return this.instanceCheck(value,DataView);
    }

    sharedArrayBuffer(value){
        return this.instanceCheck(value,(function(){try{return SharedArrayBuffer}catch{ return Fake}})());
    }

    //Intl (browser internationalization)
    intlDateTimeFormat(value){
        return this.instanceCheck(value,Intl.DateTimeFormat);
    }

    intlCollator(value){
        return this.instanceCheck(value,Intl.Collator);
    }

    intlDisplayNames(value){
        return this.instanceCheck(value,Intl.DisplayNames);
    }

    intlListFormat(value){
        return this.instanceCheck(value,Intl.ListFormat);
    }

    intlLocale(value){
        return this.instanceCheck(value,Intl.Locale);
    }

    intlNumberFormat(value){
        return this.instanceCheck(value,Intl.NumberFormat);
    }

    intlPluralRules(value){
        return this.instanceCheck(value,Intl.PluralRules);
    }

    intlRelativeTimeFormat(value){
        return this.instanceCheck(value,Intl.RelativeTimeFormat);
    }

    intlRelativeTimeFormat(value){
        return this.instanceCheck(value,Intl.RelativeTimeFormat);
    }



    //garbage collection
    finalizationRegistry(value){
        return this.instanceCheck(value,FinalizationRegistry);
    }

    weakRef(value){
        return this.instanceCheck(value,WeakRef);
    }
}

export {Is as default, Is};
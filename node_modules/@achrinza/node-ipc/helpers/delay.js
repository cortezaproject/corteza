async function delay(ms=100) {
    return new Promise(
        resolve => {
            setTimeout(resolve, ms);
        }
    );
}

export {
    delay as default,
    delay
}
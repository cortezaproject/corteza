import lcov2badge from 'lcov2badge';
import {writeFileSync} from 'fs';

lcov2badge.badge(
    './coverage/lcov.info', 
    function(err, svgBadge){
        if (err) throw err;
        writeFileSync('./coverage/lcov.svg', svgBadge); 
    }
);
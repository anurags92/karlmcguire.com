"use strict";

const render = () => {
    const map = document.getElementsByClassName("map__row");

    Array.from(map).map((row) => {
        Array.from(row.childNodes).map((box) => {
            if(box.classList == undefined) return;
           
            const random = Math.floor(Math.random() * Math.floor(3));

            if(random == 0) {
                box.classList.add("map__box--one"); 
            } else if(random == 1) {
                box.classList.add("map__box--two"); 
            } else if(random == 2) {
                box.classList.add("map__box--three"); 
            }

            box.innerHTML = `<div class="map__box__tip">April 20th, 2019</div>`;
        });
    });
};

window.onload = () => render();    

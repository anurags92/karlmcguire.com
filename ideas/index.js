window.onload = () => {
  
  const rows = Array.from(document.getElementsByClassName("row"));

  rows.map((e) => {
    const boxs = Array.from(e.childNodes);

    boxs.map((b, i) => {
      if (b.innerText != "" && b.innerText != undefined)
        b.classList.add("row__box--grow");
      else
        b.innerText = "04/18/19";
    });
  });
};

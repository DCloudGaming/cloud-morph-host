const setupModal = (modalId, btnId, spanId) => {
  var modal = document.getElementById(modalId);
  var btn = document.getElementById(btnId);
  var span = document.getElementById(spanId);
  console.log(modal);

  btn.onclick = function () {
    modal.style.display = "block";
  };

  span.onclick = function () {
    modal.style.display = "none";
  };

  window.onclick = function (event) {
    if (event.target == modal) {
      modal.style.display = "none";
    }
  };
};

setupModal("addAppModal", "addAppButton", "addAppClose");
setupModal("connectWalletModal", "connectWalletButton", "connectWalletClose");

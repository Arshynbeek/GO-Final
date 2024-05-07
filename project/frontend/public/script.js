document.addEventListener("DOMContentLoaded", function () {
  const foods = document.querySelectorAll(".food");
  const dots = document.querySelectorAll(".dot");
  let currentIndex = 0;

  function updateDots() {
    dots.forEach((dot, index) => {
      if (index === currentIndex) {
        dot.classList.add("active");
      } else {
        dot.classList.remove("active");
      }
    });
  }

  function slideCarousel() {
    const totalFoods = foods.length;
    foods.forEach((food, index) => {
      food.style.transform = `translateX(-${100 * currentIndex}%)`;
    });
    updateDots();
    currentIndex = (currentIndex + 1) % totalFoods;
  }

  dots.forEach((dot) => {
    dot.addEventListener("click", () => {
      currentIndex = parseInt(dot.getAttribute("data-index"));
      slideCarousel();
      clearInterval(autoSlide);
      autoSlide = setInterval(slideCarousel, 3000);
    });
  });

  updateDots();
  let autoSlide = setInterval(slideCarousel, 3000);


  let open = true;
  document.querySelector(".menu").addEventListener("click", () => {
    if (open) {
      document.querySelector(".menu-bar").style.display = "none";
      document.querySelector(".menu-close").style.display = "block";
      document.querySelector(".container").style.marginLeft = "400px";
      document.querySelectorAll(".menu-open").forEach(e => e.style.display = "flex");
      open = !open;
    } else {
      document.querySelector(".menu-bar").style.display = "block";
      document.querySelector(".menu-close").style.display = "none";
      document.querySelector(".container").style.marginLeft = "120px";
      document.querySelectorAll(".menu-open").forEach(e => e.style.display = "none");
      open = !open;
    }
  });
});
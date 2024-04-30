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
});

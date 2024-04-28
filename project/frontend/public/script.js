document.addEventListener("DOMContentLoaded", function () {
  const foods = document.querySelectorAll(".food");
  let currentIndex = 0;

  function slideCarousel() {
    const totalFoods = foods.length;
    foods.forEach((food, index) => {
      food.style.transform = `translateX(-${100 * currentIndex}%)`;
    });
    currentIndex = (currentIndex + 1) % totalFoods;
  }

  setInterval(slideCarousel, 3000);
});
export default (app) => {
	app.directive("debounce", {
		mounted(el, binding) {
			let timer = null;
			el.addEventListener("click", () => {
				let firstClick = !timer;

				if (firstClick) {
					binding.value();
				}
				if (timer) {
					clearTimeout(timer);
				}
				timer = setTimeout(() => {
					timer = null;
					if (!firstClick) {
						binding.value();
					}
				}, 1000);
			});
		},
	});
};

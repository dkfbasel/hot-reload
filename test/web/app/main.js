import Vue from 'vue';
import Router from 'vue-router';
import Resource from 'vue-resource';

// include the application components
import App from './app.vue';
import Start from './start.vue';

// register the router with vue
Vue.use(Router);
Vue.use(Resource);

// define router with maps and options
var router = new Router({
	history: true
});

router.map({
	'/': {
		name: 'start',
		component: Start
	}
});

router.alias({
	'*': '/'
});

router.start(App, '#app');

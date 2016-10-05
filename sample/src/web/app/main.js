import Vue from 'vue';
import Router from 'vue-router';
import Resource from 'vue-resource';

// include the application components
import App from './app.vue';
import Start from './start.vue';

// enable debugging if not in production
if (process.env.NODE_ENV !== 'production') {
	Vue.config.debug = true;
}

// register the router with vue
Vue.use(Router);
Vue.use(Resource);

const routes = [
	{
		path: '/',
		component: Start
	}
];

// define router with maps and options
var router = new Router({
	mode: 'history',
	routes: routes
});

const app = new Vue({
	router,
	el: '#app',
	render: h => h(App)
});

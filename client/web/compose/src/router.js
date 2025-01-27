import Router from 'vue-router'
import routes from './views/routes'

const router = new Router({
  mode: 'history',
  routes,
})

// Add global error handler for navigation errors
router.onError((error) => {
  console.warn('Navigation error occurred:', error)
  // Silently handle the error without crashing
})

export default router

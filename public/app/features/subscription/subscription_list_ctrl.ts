import angular from 'angular';
import { BASE_URL } from 'app/main.conf';
export class SubscriptionListCtrl {
  subscriptionType: any = ['Subscriptions', 'Open Subscritpions', 'Cancelled Subscriptions'];
  subscription: any[];
  activeSubscription: any[];
  activeTabIndex: number;
  navModel: any;
  subscriptionPlaced: string;
  searchText: string;
  constructor(private $http: any) {
    // this.navModel = navModelSrv.getSubscriptionNav();
    this.subscriptionPlaced = '';
    this.searchText = '';
    this.getSubscriptionList().then(resData => {
      this.subscription = resData.data;
      this.changeTab(0);
    });
  }

  getSubscriptionList() {
    return this.$http.get(BASE_URL + 'subscription');
  }

  changeTab(tabIndex) {
    this.activeTabIndex = tabIndex;
    this.filterSubscription(this.activeTabIndex, this.searchText);
  }

  onSearchSubscription() {
    this.filterSubscription(this.activeTabIndex, this.searchText);
  }

  filterSubscription(activeTabIndex, serchText) {
    this.activeSubscription = [];
    this.activeSubscription = this.subscription.filter(subscription => {
      let regex = new RegExp(serchText, 'ig');
      return subscription.type === activeTabIndex && regex.test(subscription.name);
    });
  }

  redirectToProjectRoom() {
    window.location.href = 'workflow';
  }
}
angular.module('grafana.controllers').controller('SubscriptionListCtrl', SubscriptionListCtrl);

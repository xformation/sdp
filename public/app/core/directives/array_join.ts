import _ from 'lodash';
import coreModule from '../core_module';
import $ from 'jquery';

export function arrayJoin() {
  'use strict';

  return {
    restrict: 'A',
    require: 'ngModel',
    link: function(scope, element, attr, ngModel) {
      function split_array(text) {
        return (text || '').split(',');
      }

      function join_array(text) {
        if (_.isArray(text)) {
          return (text || '').join(',');
        } else {
          return text;
        }
      }

      ngModel.$parsers.push(split_array);
      ngModel.$formatters.push(join_array);
    },
  };
}

coreModule.directive('arrayJoin', arrayJoin);

function drawChart() {
  return {
    restrict: 'A',
    scope: {
      chartPoints: '=',
      type: '@',
      showLegends: '@',
    },
    controllerAs: 'ctrl',
    link: function(scope, element) {
      let Chart = window['Chart'];
      let $canvas = $(element).find('canvas');
      if (scope.chartPoints) {
        let datasets = JSON.parse(scope.chartPoints);
        let chart = new Chart($canvas[0], {
          type: scope.type,
          data: {
            // labels:['a','b','c'],
            datasets: datasets,
          },
          options: {
            scales: {
              yAxes: [
                {
                  ticks: {
                    beginAtZero: true,
                  },
                },
              ],
            },
            legend: {
              position: 'right',
              display: scope.showLegends === 'true',
            },
          },
        });
        console.log(chart);
      }
    },
  };
}

coreModule.directive('drawchart', drawChart);

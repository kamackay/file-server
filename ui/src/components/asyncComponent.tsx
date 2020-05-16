import * as React from "react";

interface AsyncComponentState {
  component: any;
}

export default (importComponent: any) => {
  class AsyncComponent extends React.Component<any, AsyncComponentState> {
    constructor(props: any) {
      super(props);

      this.state = {
        component: null
      };
    }

    public async componentDidMount() {
      const { default: component } = await importComponent();

      this.setState({ component });
    }

    public render() {
      const C = this.state.component;

      return C ? <C {...this.props} /> : null;
    }
  }

  return AsyncComponent;
};

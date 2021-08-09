import React from 'react';
import clsx from 'clsx';
import styles from './HomepageFeatures.module.css';

const FeatureList = [
  {
    title: 'A Solution to store Private keys and Sign txns securely',
    Svg: require('../../static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        This Plugin is design in a way that stores a user's mnemonic inside vault in an encrypted manner.
      </>
    ),
  },
  {
    title: 'Easy User Management',
    Svg: require('../../static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
          Single unique id(UUID) of the user which will be used to access the user's keys stored in the vault.
      </>
    ),
  },
    {
        title: 'Focus on What Matters',
        Svg: require('../../static/img/undraw_docusaurus_mountain.svg').default,
        description: (
            <>
                This plugin help you to focus on other development rather focusing on signing a txns manually.
            </>
        ),
    },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} alt={title} />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

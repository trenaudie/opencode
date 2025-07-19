import { makeScene2D, Rect, Circle } from '@motion-canvas/2d';
import { all, createRef } from '@motion-canvas/core';

import { createSignal, createEffect } from '@motion-canvas/core';

export default makeScene2D(function* (view) {
  view.fill('#000000');
  const square = createRef<Circle>();

  // Signals for position
  const radius = 200;
  const numFrames = 120;
  const angle = createSignal(0);
  const x = createSignal(0);
  const y = createSignal(0);

  view.add(
    <Circle
      ref={square}
      width={140} height={140}
      fill="#0000FF"
    />
  );

  // Computing the position using effect
  createEffect(() => {
    const currentAngle = angle();
    x(radius * Math.cos(currentAngle));
    y(radius * Math.sin(currentAngle));
  });

  // Animate the circle to move in a circular path
  for (let frame = 0; frame <= numFrames; frame++) {
    yield* all(
      angle((Math.PI * 2 * frame) / numFrames, 0.1)
    );
    yield* square().position(x(), y(), 1);
  }
});
  const square = createRef<Circle>();
  view.fill('#000000');
  view.add(
    <Circle
      ref={square}
      x={0}
      width={140} height={140}
      fill="#0000FF"
    />
  );

  // Define the circle's movement parameters
  const radius = 200;
  const numFrames = 120;
  let angle = 0;

  // Animate the circle to move in a circular path
  while (angle <= Math.PI * 2) {
    const x = radius * Math.cos(angle);
    const y = radius * Math.sin(angle);
    yield* square().position(x, y, 1);
    angle += Math.PI * 2 / numFrames;
  }
});
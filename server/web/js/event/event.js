/**
 * Event publishing / subscribe module.
 * Just a simple observer pattern.
 * @version 1
 */
const event = (() => {
  const topics = {};

  return {
    /**
     * Subscribes onto some event.
     *
     * @param topic The name of the event.
     * @param listener A callback function to call during the event.
     * @param order A number in a queue of event handlers to run callback in ordered manner.
     * @returns {{unsub: unsub}} The function to remove this subscription.
     * @example
     * const sub01 = event.sub('rapture', () => {a}, 1)
     * ...
     * sub01.unsub()
     */
    sub: (topic, listener, order) => {
      if (!topics[topic]) topics[topic] = [];

      // order handling stuff
      const value = { order: order || 0, listener: listener };
      topics[topic].push(value);
      topics[topic].sort((a, b) => a.order - b.order);
      const index = topics[topic].indexOf(value);

      return {
        unsub: () => {
          topics[topic].splice(index, 1);
        },
      };
    },

    /**
     * Publishes some event for handling.
     *
     * @param topic The name of the event.
     * @param data Additional data for the event handling.
     * Because of compatibility we have to use a dumb obj wrapper {a: a, b: b} for params instead of (topic, ...data).
     * @example
     * event.pub('rapture', {time: now()})
     */
    pub: (topic, data) => {
      if (!topics[topic] || topics[topic].length < 1) return;

      topics[topic].forEach((listener) => {
        listener.listener(data !== undefined ? data : {});
      });
    },
  };
})();

// events
// WebRTC sign
const MEDIA_STREAM_INITIALIZED = "mediaStreamInitialized";
const MEDIA_STREAM_SDP_AVAILABLE = "mediaStreamSdpAvailable";
const MEDIA_STREAM_CANDIDATE_ADD = "mediaStreamCandidateAdd";
const MEDIA_STREAM_CANDIDATE_FLUSH = "mediaStreamCandidateFlush";
const MEDIA_STREAM_READY = "mediaStreamReady";
const PING_REQUEST = "pingRequest";
const PING_RESPONSE = "pingResponse";

const CONNECTION_READY = "connectionReady";
const CONNECTION_CLOSED = "connectionClosed";

const KEY_PRESSED = "keyPressed";
const KEY_RELEASED = "keyReleased";
const KEY_STATE_UPDATED = "keyStateUpdated";
const KEYBOARD_KEY_PRESSED = "keyboardKeyPressed";

const MOUSE_DOWN = "mouseDown";
const MOUSE_UP = "mouseUp";
const MOUSE_MOVE = "mouseMove";

const CLIENT_INIT = "clientInit";
const HOST_AVAILABLE = "hostAvailable";
